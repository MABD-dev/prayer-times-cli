package ui

import (
	"os"

	"github.com/aquasecurity/table"
	"github.com/fatih/color"
	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

var (
	prayerTimeHeaderrFgColor = color.New(color.FgHiGreen)
)

func RenderPrayerTime(dayPrayerTimes models.DayPrayerTimes) {
	table := table.New(os.Stdout)

	table.SetHeaders(
		prayerTimeHeaderrFgColor.Sprint("Fajr"),
		prayerTimeHeaderrFgColor.Sprint("Dhuhr"),
		prayerTimeHeaderrFgColor.Sprint("Asr"),
		prayerTimeHeaderrFgColor.Sprint("Maghrib"),
		prayerTimeHeaderrFgColor.Sprint("Isha"),
	)
	times := dayPrayerTimes.PrayerTimes
	table.AddRow(
		times.Fajr,
		times.Dhuhr,
		times.Asr,
		times.Maghrib,
		times.Isha,
	)
	table.Render()
}
