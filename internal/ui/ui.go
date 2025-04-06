package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/aquasecurity/table"
	"github.com/fatih/color"
	"github.com/mabd-dev/prayer-times-cli/internal/domain"
)

var (
	prayerTimeHeaderrFgColor = color.New(color.FgHiGreen)
	remainingTimeFgColor     = color.New(color.FgHiGreen)
)

func RenderPrayerTime(prayers []domain.Prayer) {
	table := table.New(os.Stdout)

	headers := []string{}
	prayerTimes := []string{}
	for _, p := range prayers {
		headers = append(headers, prayerTimeHeaderrFgColor.Sprint(p.Name))

		timeFormatted := p.Time.Format("3:04 pm")
		prayerTimes = append(prayerTimes, timeFormatted)
	}
	table.SetHeaders(headers...)
	table.AddRow(prayerTimes...)
	table.Render()
}

func RenderDate(time time.Time) {
	formatted := time.Format("Monday 02/01/2006")
	fmt.Println(formatted)
}

func RenderTimeRemaining(
	prayer domain.Prayer,
	duration time.Duration,
) {
	coloredPrayerName := remainingTimeFgColor.Sprint(prayer.Name)
	coloredHours := remainingTimeFgColor.Sprint(int(duration.Hours()))
	coloredMinutes := remainingTimeFgColor.Sprint(int(duration.Minutes()) % 60)
	fmt.Printf("%v hours, %v minutes to %v\n", coloredHours, coloredMinutes, coloredPrayerName)
}
