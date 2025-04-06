package domain

import (
	"fmt"
	"time"
)

// GetDayPrayerTimeFor get caches data locally or fetch new data from remote then save locally.
// Then search data for specific @year @month and @day. If found return prayer times
func GetDayPrayerTimeFor(time time.Time) *DayPrayers {
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

	prayerTimes := getPrayerTimes(*data, dateStr)
	return mapToDayPrayer(*prayerTimes)
}

// getNextAndPreviousPrayerTimes
//
// @returns
//   - (previous prayer, next prayer)
func GetNextAndPreviousPrayerTimes(dayPrayers DayPrayers) (*Prayer, *Prayer) {
	day := time.Now()

	yesterdayDate := day.Add(-24 * time.Hour)
	yesterdayPrayers := GetDayPrayerTimeFor(yesterdayDate)
	if yesterdayPrayers == nil {
		return nil, nil //errors.New("Failed to get yesterday's prayers")
	}
	fmt.Println(yesterdayPrayers)

	tomorrowDate := day.Add(24 * time.Hour)
	tomorrowPrayers := GetDayPrayerTimeFor(tomorrowDate)
	if tomorrowPrayers == nil {
		return nil, nil //errors.New("Failed to get tomorrow's prayers")
	}
	fmt.Println(tomorrowPrayers)

	combinedPrayerTimes := []Prayer{}

	// take only fajr prayer from next day. best case it would be previous prayer
	combinedPrayerTimes = append(combinedPrayerTimes, (*yesterdayPrayers).Prayers[len((*yesterdayPrayers).Prayers)-1])

	combinedPrayerTimes = append(combinedPrayerTimes, dayPrayers.Prayers...)

	// take only fajr prayer from next day. best case it would be next prayer
	combinedPrayerTimes = append(combinedPrayerTimes, (*tomorrowPrayers).Prayers[0])

	nextPrayerIndex := -1
	for i, p := range combinedPrayerTimes {
		if p.Time.After(day) || p.Time.Equal(day) {
			nextPrayerIndex = i
			break
		}
	}

	return &combinedPrayerTimes[nextPrayerIndex-1], &combinedPrayerTimes[nextPrayerIndex]
}

// GetTimeRemainingTo
//
// @Returns
//   - hours remaining
//   - minutes remaining
func GetTimeRemainingTo(nextPrayerTime time.Time) *time.Duration {
	now := time.Now()
	if now.After(nextPrayerTime) {
		return nil
	}

	duration := nextPrayerTime.Sub(now)
	return &duration

}
