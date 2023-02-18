package core

import (
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/kelvins/sunrisesunset"
	"io"
	"strconv"
	"time"
)

/*
What I want to display: percent through the apparent solar day we are
*/

//go:embed cities.csv
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

func (c *CityInfo) GetSunriseSunset(at time.Time) (sunrise time.Time, sunset time.Time, err error) {
	_, secondsOffset := at.Zone()
	hoursOffset := float64(secondsOffset) / 60 / 60
	p := sunrisesunset.Parameters{
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
		UtcOffset: hoursOffset,
		Date:      at,
	}

	return p.GetSunriseSunset()
}

func GetPeriodPercent(c *CityInfo, at time.Time, debug bool) (string, error) {
	if debug {
		fmt.Printf("City: %#v\n", c)
	}

	sunrise, sunset, err := c.GetSunriseSunset(at)
	if err != nil {
		return "", err
	}

	if debug {
		fmt.Println("Sunrise:", sunrise.Format("15:04:05")) // Sunrise: 06:11:44
		fmt.Println("Sunset:", sunset.Format("15:04:05"))   // Sunset: 18:14:27
		fmt.Println("Length of apparent solar time in mean solar time:", sunset.Sub(sunrise))
		fmt.Println("Solar noon:", sunrise.Add(sunset.Sub(sunrise)/2).Format("15:04:05"))
	}

	dayDuration := sunset.Sub(sunrise)
	if at.After(sunrise) && at.Before(sunset) {
		hasPassed := at.Sub(sunrise)
		if debug {
			fmt.Println("Time passed since sunrise:", hasPassed)
		}
		fractionPassed := hasPassed.Seconds() / dayDuration.Seconds()
		return fmt.Sprintf("%.0f%% ☉", fractionPassed*100), nil
	}

	// https://glossary.ametsoc.org/wiki/Mean_solar_day
	meanDayDuration, err := time.ParseDuration("86400s")
	if err != nil {
		return "", err
	}
	nightDuration := meanDayDuration - dayDuration

	if at.Before(sunrise) {
		nightEnd := sunrise
		nightStart := nightEnd.Add(-nightDuration)
		hasPassedNight := at.Sub(nightStart)
		fractionPassedNight := hasPassedNight.Seconds() / nightDuration.Seconds()
		if debug {
			fmt.Println("Time passed since sunset:", hasPassedNight)
		}
		return fmt.Sprintf("%.0f%% ☾", fractionPassedNight*100), nil
	}

	if at.After(sunset) {
		nightStart := sunset
		hasPassedNight := at.Sub(nightStart)
		fractionPassedNight := hasPassedNight.Seconds() / nightDuration.Seconds()
		if debug {
			fmt.Println("Time passed since sunset:", hasPassedNight)
		}
		return fmt.Sprintf("%.0f%% ☾", fractionPassedNight*100), nil
	}

	panic("This should never happen")
}

type NarrowingError struct {
	Cities []*CityInfo
}

func (e *NarrowingError) Error() string {
	return "could not narrow between cities"
}

func FindCity(name, countryCode, fipsCode string) (*CityInfo, error) {
	// TODO cache successful runs and don't bother with opening this
	file, err := fs.Open("cities.csv")
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
		return nil, &NarrowingError{Cities: cities}
	}
	return cities[0], nil
}
