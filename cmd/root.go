/*
Copyright © 2025 MABD-dev <mabd.universe@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

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

		isToday := domain.SameDay(now, requestedDate)
		if isToday {
			fmt.Println("Calling some other function")
		} else {
			dailyPrayerSchedule, err := domain.GetDailyPrayerSchedule(requestedDate)
			if err != nil {
				return err
			}
			ui.RenderDailyPrayerSchedule(dailyPrayerSchedule)
		}
		// ui.RenderPrayerTime((*dayPrayer).Prayers)
		//
		// if isToday {
		// 	previousPrayer, nextPrayer := domain.GetNextAndPreviousPrayerTimes(*dayPrayer)
		// 	if previousPrayer == nil || nextPrayer == nil {
		// 		return errors.New("Could not get previous or next prayer")
		// 	}
		//
		// 	reminaingToNextPrayer := domain.GetTimeRemainingTo(nextPrayer.Time)
		// 	if reminaingToNextPrayer == nil {
		// 		return errors.New("Failed to get time remaining to next prayer")
		// 	}
		// 	ui.RenderTimeRemaining(*nextPrayer, *reminaingToNextPrayer)
		//
		// 	timeProgressPercent := domain.TimeProgressPercent(previousPrayer.Time, nextPrayer.Time)
		// 	ui.RenderTimeProgress(*previousPrayer, *nextPrayer, timeProgressPercent)
		// }

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
