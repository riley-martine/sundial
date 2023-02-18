package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/riley-martine/sundial/internal/core"
	"github.com/spf13/cobra"
)

// Set by goreleaser:
// https://goreleaser.com/cookbooks/using-main.version/?h=version
var version = "dev"

var (
	debug       bool
	cityName    string
	countryCode string
	fipsCode    string
	givenTime   string
)

var rootCmd = &cobra.Command{
	Use:     "sundial --city CITY",
	Short:   "Print the percent through the day or night.",
	Version: version,
	Long: `Sundial is a program to print the percent through the day or night.
https://github.com/riley-martine/sundial`,
	Run: func(cmd *cobra.Command, args []string) {
		city, err := core.FindCity(cityName, countryCode, fipsCode)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if strings.HasPrefix(err.Error(), "could not narrow") {
				fmt.Fprintln(os.Stderr, "You may need to be more specific about which city you're in. Try specifying a country code (second field) and a fips code (third field).")
				fmt.Fprintln(os.Stderr, "    e.g. sundial -city Washington -country US -fipscode IL")
			}
			os.Exit(1)
		}

		// This probably doesn't need ParseInLocation
		// Adding that data back to the cities CSV balloons the size
		t := time.Now()
		if givenTime != "" {
			t, err = time.Parse(time.UnixDate, givenTime)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		out, err := core.GetPeriodPercent(city, t, debug)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(out)
	},
}

func init() {
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Print debug logging. Default: false")

	rootCmd.Flags().StringVar(&cityName, "city", "", "Name of city you're in. Required.")
	rootCmd.MarkFlagRequired("city")
	rootCmd.Flags().StringVar(&countryCode, "country", "", "Two-letter country code, e.g. 'US'. Not required if only one city with name.")
	rootCmd.Flags().StringVar(&fipsCode, "fipscode", "", `Fipscode of region you're in. In the US, this is the two-letter state abbreviation.
Otherwise, search http://download.geonames.org/export/dump/admin1CodesASCII.txt
for '$countryCode.' and select the value after the period for the region you're in.
Not required if only one city in country with name.`)

	rootCmd.Flags().StringVar(&givenTime, "time", "", "Time to convert, in time.UnixDate format. Defaults to now.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}