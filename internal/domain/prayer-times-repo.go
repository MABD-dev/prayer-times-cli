package domain

import (
	"errors"
	"fmt"
	"time"
)

func GetDailyPrayerSchedule(date time.Time) (DailyPrayerSchedule, error) {
	dayPrayers := GetDayPrayerTimeFor(date)
	if dayPrayers == nil {
		return DailyPrayerSchedule{}, errors.New("Failed to get day prayer")
	}

	return DailyPrayerSchedule{
		Date:    dayPrayers.Date,
		Prayers: dayPrayers.Prayers,
	}, nil
}

func GetActivePrayerTracking(date time.Time) (ActivePrayerTracking, error) {
	dayPrayers := GetDayPrayerTimeFor(date)
	if dayPrayers == nil {
		return ActivePrayerTracking{}, errors.New("Failed to get day prayer")
	}

	previousPrayer, nextPrayer := GetNextAndPreviousPrayerTimes(*dayPrayers)
	if previousPrayer == nil || nextPrayer == nil {
		return ActivePrayerTracking{}, errors.New("Could not get previous or next prayer")
	}

	reminaingToNextPrayer := GetTimeRemainingTo(nextPrayer.Time)
	if reminaingToNextPrayer == nil {
		return ActivePrayerTracking{}, errors.New("Failed to get time remaining to next prayer")
	}

	timeProgressPercent := TimeProgressPercent(previousPrayer.Time, nextPrayer.Time)

	return ActivePrayerTracking{
		DailyPrayerSchedule: DailyPrayerSchedule{
			Date:    dayPrayers.Date,
			Prayers: dayPrayers.Prayers,
		},
		NextPrayer:    nextPrayer.Name,
		TimeRemaining: *reminaingToNextPrayer,
		Progress:      timeProgressPercent,
	}, nil
}

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
	if prayerTimes == nil {
		return nil
	}
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

	tomorrowDate := day.Add(24 * time.Hour)
	tomorrowPrayers := GetDayPrayerTimeFor(tomorrowDate)
	if tomorrowPrayers == nil {
		return nil, nil //errors.New("Failed to get tomorrow's prayers")
	}

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

func TimeProgressPercent(
	previousPrayerTime time.Time,
	nextPrayerTime time.Time,
) float64 {
	now := time.Now()
	totalDuration := nextPrayerTime.Sub(previousPrayerTime).Seconds()
	passedDuration := nextPrayerTime.Sub(now).Seconds()

	if passedDuration <= 0 {
		return 100.0
	}

	percent := 100 - (passedDuration/totalDuration)*100.0
	if percent < 0 {
		return 0.0
	}
	if percent > 100.0 {
		return 100.0
	}
	return percent
}
