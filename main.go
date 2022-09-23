package main

import (
	"fmt"
	"time"

	"github.com/kelvins/sunrisesunset"
)

/*
What I want to display: percent through the apparent solar day we are
*/

func main() {
	// You can use the Parameters structure to set the parameters
	now := time.Now()
	p := sunrisesunset.Parameters{
		Latitude:  39.76, // Denver, CO
		Longitude: -104.9,
		UtcOffset: -6.0, // MDT
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
