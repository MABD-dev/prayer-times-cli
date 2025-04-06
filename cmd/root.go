/*
Copyright © 2025 MABD-dev <mabd.universe@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mabd-dev/prayer-times-cli/internal/domain"
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

		now := time.Now()
		requestedDate := time.Date(year, time.Month(month), day, now.Hour(), now.Minute(), 0, 0, now.Location())

		prayerTime := domain.GetDayPrayerTimeFor(requestedDate)
		if prayerTime == nil {
			fmt.Println("could not find prayer time")
			return nil
		}
		ui.RenderPrayerTime(*prayerTime)

		nextPrayerTime, prayerName := domain.GetNextPrayerTime(requestedDate, *prayerTime)
		if nextPrayerTime == nil {
			return fmt.Errorf("Could not next prayer time!")
		}

		timeToNextPrayer := nextPrayerTime.Sub(now)

		green := color.New(color.FgHiGreen)
		coloredPrayerName := green.Sprint(prayerName)
		coloredHours := green.Sprint(int(timeToNextPrayer.Hours()))
		coloredMinutes := green.Sprint(int(timeToNextPrayer.Minutes()) % 60)
		fmt.Printf("%v hours, %v minutes till %v\n", coloredHours, coloredMinutes, coloredPrayerName)

		// previousPrayerName, err := domain.GetPreviousPrayerName(prayerName)
		// fmt.Printf("%v to %v\n", previousPrayerName, prayerName)

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
