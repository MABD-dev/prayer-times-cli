package domain

import (
	"fmt"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

// GetDayPrayerTimeFor get caches data locally or fetch new data from remote then save locally.
// Then search data for specific @year @month and @day. If found return prayer times
func GetDayPrayerTimeFor(time time.Time) *models.DayPrayerTimes {
	dateStr := formatDate(time)
	localYearFilename := fmt.Sprintf("%v.json", time.Year())

	data := loadFromLocal(localYearFilename)
	if data == nil {
		res, err := fetchAndSavePrayerTimes(time.Year(), localYearFilename)
		if err != nil {
			fmt.Println("Failed to fetch data from internet")
			return nil
		}
		data = res
	}

	return getPrayerTimes(*data, dateStr)
}

func GetNextPrayerTime(
	requestedTime time.Time,
	dayPrayerTimes models.DayPrayerTimes,
) *time.Time {
	nextPrayerTime := getNextPrayerTime(requestedTime, requestedTime, dayPrayerTimes.PrayerTimes)
	if nextPrayerTime == nil {
		nextDay := requestedTime.Add(24 * time.Hour)
		nextDayPrayerTimes := GetDayPrayerTimeFor(nextDay)
		if nextDayPrayerTimes == nil {
			return nil
		}
		nextPrayerTime = getNextPrayerTime(requestedTime, nextDay, dayPrayerTimes.PrayerTimes)
	}
	return nextPrayerTime
}
