package main

import (
	"embed"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kelvins/sunrisesunset"
)

/*
What I want to display: percent through the apparent solar day we are
*/

var cityName = flag.String("city", "", "Name of city you're in. Required.")
var countryCode = flag.String("country", "", "Two-letter country code, e.g. 'US'. Not required if only one city by name.")
var fipsCode = flag.String("fipscode", "", "Fipscode of region you're in. In the US, this is the two-letter state abbreviation. Otherwise, search http://download.geonames.org/export/dump/admin1CodesASCII.txt for '$countryCode.' and select the value after the period for the region you're in. Not required if only one city in country with name.")

//go:embed static/cities.csv
var fs embed.FS

type CityInfo struct {
	Name        string
	CountryCode string
	FipsCode    string
	Latitude    float64
	Longitude   float64
}

func (c *CityInfo) String() string {
	return fmt.Sprintf("{%s, %s, %s, %0.2f, %0.2f}", c.Name, c.CountryCode, c.FipsCode, c.Latitude, c.Longitude)
}

func findCity(name, countryCode, fipsCode string) (*CityInfo, error) {
	// TODO cache successful runs and don't bother with opening this
	file, err := fs.Open("static/cities.csv")
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	r.Comma = '\t'
	r.ReuseRecord = true

	cities := []*CityInfo{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if record[0] != name {
			continue
		}

		if countryCode != "" && record[1] != countryCode {
			continue
		}

		if fipsCode != "" && record[2] != fipsCode {
			continue
		}
		lat, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		long, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		cities = append(cities, &CityInfo{Name: record[0], CountryCode: record[1], FipsCode: record[2], Latitude: lat, Longitude: long})
	}
	if len(cities) == 0 {
		return nil, fmt.Errorf("unable to find city '%s' in country '%s' with fips code '%s'", name, countryCode, fipsCode)
	}
	if len(cities) > 1 {
		return nil, fmt.Errorf("could not narrow down between cities: %+v", cities)
	}
	return cities[0], nil
}

func main() {
	flag.Parse()
	if *cityName == "" {
		fmt.Println("City is a required argument. Usage: sundial -city Denver")
		os.Exit(1)
	}

	city, err := findCity(*cityName, *countryCode, *fipsCode)
	if err != nil {
		fmt.Println(err)
		if strings.HasPrefix(err.Error(), "could not narrow") {
			fmt.Println("You may need to be more specific about which city you're in. Try specifying a country code (second field) and a fips code (third field).")
			fmt.Println("e.g. sundial -city Washington -country US -fipscode IL")
		}
		os.Exit(1)
	}

	now := time.Now()
	_, secondsOffset := now.Zone()
	hoursOffset := float64(secondsOffset) / 60 / 60
	p := sunrisesunset.Parameters{
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		UtcOffset: hoursOffset,
		Date:      now,
	}

	// Calculate the sunrise and sunset times
	sunrise, sunset, err := p.GetSunriseSunset()

	if err != nil {
		panic(err)
	}

	// fmt.Println("Sunrise:", sunrise.Format("15:04:05")) // Sunrise: 06:11:44
	// fmt.Println("Sunset:", sunset.Format("15:04:05"))   // Sunset: 18:14:27
	// fmt.Println("Length of apparent solar time in mean solar time:", sunset.Sub(sunrise))
	//fmt.Println("Solar noon:", sunrise.Add(sunset.Sub(sunrise)/2).Format("15:04:05"))
	dayDuration := sunset.Sub(sunrise)
	if now.After(sunrise) && now.Before(sunset) {
		hasPassed := now.Sub(sunrise)
		// fmt.Println("Time passed since sunrise:", hasPassed)
		fractionPassed := hasPassed.Seconds() / dayDuration.Seconds()
		fmt.Printf("%.0f%% ☉", fractionPassed*100)
		return
	}

	// https://glossary.ametsoc.org/wiki/Mean_solar_day
	meanDayDuration, err := time.ParseDuration("86400s")
	if err != nil {
		panic(err)
	}
	nightDuration := meanDayDuration - dayDuration

	if now.Before(sunrise) {
		nightEnd := sunrise
		nightStart := nightEnd.Add(-nightDuration)
		hasPassedNight := now.Sub(nightStart)
		fractionPassedNight := hasPassedNight.Seconds() / nightDuration.Seconds()

		fmt.Printf("%.0f%% ☾", fractionPassedNight*100)
		return
	}

	if now.After(sunset) {
		nightStart := sunset
		hasPassedNight := now.Sub(nightStart)
		fractionPassedNight := hasPassedNight.Seconds() / nightDuration.Seconds()

		fmt.Printf("%.0f%% ☾", fractionPassedNight*100)
		return
	}

	panic("This should never happen")
}
