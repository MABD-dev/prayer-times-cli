package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aquasecurity/table"
	"github.com/fatih/color"
	"github.com/mabd-dev/prayer-times-cli/internal/domain"
)

var (
	prayerTimeHeaderrFgColor = color.New(color.FgHiGreen)
	remainingTimeFgColor     = color.New(color.FgHiGreen)
	timeProgressFgColor      = color.New(color.FgHiGreen)
	timeProgressBgColor      = color.New(color.BgHiGreen)
)

func RenderDailyPrayerSchedule(dailyPrayerSchedule domain.DailyPrayerSchedule) {
	RenderDate(dailyPrayerSchedule.Date)
	RenderPrayerTimes(dailyPrayerSchedule.Prayers)
}

func RenderActivePrayerTracking(activePrayerTracking domain.ActivePrayerTracking) {
	RenderDate(activePrayerTracking.Date)
	RenderPrayerTimes(activePrayerTracking.Prayers)
	RenderTimeRemaining(activePrayerTracking.NextPrayer, activePrayerTracking.TimeRemaining)
	RenderTimeProgress(
		activePrayerTracking.PreviousPrayer,
		activePrayerTracking.NextPrayer,
		activePrayerTracking.Progress,
	)
}

func RenderPrayerTimes(prayers []domain.Prayer) {
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

// RenderDate format date and draw it on screen
func RenderDate(time time.Time) {
	formatted := time.Format("Monday 02/01/2006")
	fmt.Println(formatted)
}

// RenderTimeRemaining show how many hours and minutes remaining till next prayer +
// show next prayer name
func RenderTimeRemaining(
	nextPrayer string,
	duration time.Duration,
) {
	coloredPrayerName := remainingTimeFgColor.Sprint(nextPrayer)
	coloredHours := remainingTimeFgColor.Sprint(int(duration.Hours()))
	coloredMinutes := remainingTimeFgColor.Sprint(int(duration.Minutes()) % 60)
	fmt.Printf("%v hours, %v minutes to %v\n", coloredHours, coloredMinutes, coloredPrayerName)
}

// RenderTimeProgress shows previous and next prayer names and in between progress bar like
func RenderTimeProgress(
	previousPrayer string,
	nextPrayer string,
	timeProgressPercent float64,
) {
	symbol := "─"
	totalNumberOfSymbols := 40
	coloredSymbols := int(timeProgressPercent / 100 * float64(totalNumberOfSymbols))
	whiteSymbols := totalNumberOfSymbols - coloredSymbols

	var sb strings.Builder

	sb.WriteString(previousPrayer)
	sb.WriteString(" ")
	for range coloredSymbols {
		sb.WriteString(timeProgressBgColor.Sprint(" "))
	}
	for range whiteSymbols {
		sb.WriteString(symbol)
	}
	sb.WriteString(" ")
	sb.WriteString(nextPrayer)
	fmt.Println(sb.String())

}
