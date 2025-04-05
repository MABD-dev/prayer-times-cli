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

		prayerTime := domain.GetDayPrayerTimeFor(requestedDate)
		if prayerTime == nil {
			fmt.Println("could not find prayer time")
			return nil
		}
		fmt.Printf("date=%v\n", prayerTime.Gregorian)
		ui.RenderPrayerTime(*prayerTime)

		nextPrayerTime := domain.GetNextPrayerTime(requestedDate, *prayerTime)
		if nextPrayerTime != nil {
			fmt.Printf("next prayer time=%v\n", nextPrayerTime)
		}

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
