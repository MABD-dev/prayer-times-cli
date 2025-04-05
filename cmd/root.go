/*
Copyright © 2025 MABD-dev <mabd.universe@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/domain"
	"github.com/mabd-dev/prayer-times-cli/internal/models"
	"github.com/mabd-dev/prayer-times-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prayers",
	Short: "Get prayer times for today",
	RunE: func(cmd *cobra.Command, args []string) error {

		year, err := cmd.Flags().GetInt("year")
		if err != nil {
			panic(err)
		}
		month, err := cmd.Flags().GetInt("month")
		if err != nil {
			return err
		}
		day, err := cmd.Flags().GetInt("day")
		if err != nil {
			return err
		}

		prayerTime := domain.GetDayPrayerTimeFor(year, month, day)
		if prayerTime == nil {
			fmt.Println("could not find prayer time")
			return nil
		}
		ui.RenderPrayerTime(*prayerTime)

		// now := time.Now()
		// fmt.Printf("currenet time: hour=%v, minute=%v\n", now.Hour(), now.Minute())
		// hour, minute, err := prayerTimeToHourMinute(prayerTime.PrayerTimes.Isha)
		// fmt.Printf("Isha prayer: hour=%v, minute=%v, err=%v\n", hour, minute, err)

		nextPrayerTime := getNextPrayerTime((*prayerTime).PrayerTimes)
		if nextPrayerTime == nil {
			fmt.Println("could not find next prrayer time")
		}
		fmt.Printf("next prayer time=%v\n", nextPrayerTime)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand()

	now := time.Now()
	rootCmd.PersistentFlags().IntP("year", "y", now.Year(), "Set year")
	rootCmd.PersistentFlags().IntP("month", "m", int(now.Month()), "Set month")
	rootCmd.PersistentFlags().IntP("day", "d", now.Day(), "Set day")
}

// prayerTimeTo24format takes prayer time in this format "hh:mm am"
// and return hour (in 24-hour format) and minute
func prayerTimeToHourMinute(value string) (int, int, error) {
	value = strings.ToLower(strings.TrimSpace(value))

	t, err := time.Parse("3:04 pm", value)
	if err != nil {
		return 0, 0, nil
	}
	return t.Hour(), t.Minute(), nil

	// a := strings.Split(value, " ")
	// ampm := strings.ToLower(a[1])
	//
	// time := strings.Split(a[0], ":")
	// hour, err := strconv.Atoi(time[0])
	// if err != nil {
	// 	return -1, -1
	// }
	//
	// minute, err := strconv.Atoi(time[1])
	// if err != nil {
	// 	return -1, -1
	// }
	//
	// if ampm == "pm" {
	// 	hour += 12
	// }
	//
	// return hour, minute
}

// ParseTime takes a time string like this "12:05 pm" and convert it to @time.Time
// or returns error if failed to parse
func parseTime(tStr string) (time.Time, error) {
	layout := "3:04 pm"
	now := time.Now()
	t, err := time.Parse(layout, strings.ToLower(strings.TrimSpace(tStr)))
	if err != nil {
		return time.Time{}, err
	}
	// Attach today's date to the parsed time
	todayTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 59, 0, now.Location())
	return todayTime, nil
}

// getNextPrayerTime for now assuming their is a next prayer in the day
// TODO: write better docs
func getNextPrayerTime(prayerTimes models.PrayerTimes) *time.Time {
	now := time.Now()
	//now = time.Date(now.Year(), now.Month(), now.Day(), 12, 40, 1, 0, now.Location())

	sortedPrayerTimes := []string{
		prayerTimes.Fajr,
		prayerTimes.Dhuhr,
		prayerTimes.Asr,
		prayerTimes.Maghrib,
		prayerTimes.Isha,
	}

	for _, prayerTime := range sortedPrayerTimes {
		t, err := parseTime(prayerTime)
		if err != nil {
			return nil
		}

		if t.After(now) || t.Equal(now) {
			return &t
		}
	}

	return nil
}
