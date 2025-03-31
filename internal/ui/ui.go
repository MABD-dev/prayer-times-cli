package ui

import (
	"fmt"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

func RenderPrayerTime(dayPrayerTimes models.DayPrayerTimes) {
	fmt.Println("***********")
	fmt.Printf("%v\n", dayPrayerTimes.Gregorian)
	fmt.Printf("Fajr=%v\n", dayPrayerTimes.PrayerTimes.Fajr)
	fmt.Printf("Dhuhr=%v\n", dayPrayerTimes.PrayerTimes.Dhuhr)
	fmt.Printf("Asr=%v\n", dayPrayerTimes.PrayerTimes.Asr)
	fmt.Printf("Maghrib=%v\n", dayPrayerTimes.PrayerTimes.Maghrib)
	fmt.Printf("Isha=%v\n", dayPrayerTimes.PrayerTimes.Isha)
	fmt.Println("***********")
}
