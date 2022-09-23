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

	// If no error has occurred, print the results
	if err == nil {
		// fmt.Println("Sunrise:", sunrise.Format("15:04:05")) // Sunrise: 06:11:44
		// fmt.Println("Sunset:", sunset.Format("15:04:05"))   // Sunset: 18:14:27
		// fmt.Println("Length of apparent solar time in mean solar time:", sunset.Sub(sunrise))
		//fmt.Println("Solar noon:", sunrise.Add(sunset.Sub(sunrise)/2).Format("15:04:05"))

		if now.Before(sunrise) {
			fmt.Println("0%")
		} else if now.After(sunset) {
			fmt.Println("100%")
		} else { // After sunrise and before sunset
			dayDuration := sunset.Sub(sunrise)
			hasPassed := now.Sub(sunrise)
			// fmt.Println("Time passed since sunrise:", hasPassed)
			fractionPassed := hasPassed.Seconds() / dayDuration.Seconds()
			fmt.Printf("%.0f%%", fractionPassed*100)
		}
	} else {
		fmt.Println(err)
	}
}
