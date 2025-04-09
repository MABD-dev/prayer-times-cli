package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

type PrayerTimesRepo interface {
	GetDailyPrayerSchedule(date time.Time) (DailyPrayerSchedule, error)
	GetActivePrayerTracking(date time.Time) (ActivePrayerTracking, error)
}

type PrayerTimesRepoImpl struct{}

func (r *PrayerTimesRepoImpl) GetDailyPrayerSchedule(date time.Time) (DailyPrayerSchedule, error) {
	dayPrayers := getDayPrayerTimeFor(date)
	if dayPrayers == nil {
		return DailyPrayerSchedule{}, errors.New("Failed to get day prayer")
	}

	return DailyPrayerSchedule{
		Date:    dayPrayers.Date,
		Prayers: dayPrayers.Prayers,
	}, nil
}

func (r *PrayerTimesRepoImpl) GetActivePrayerTracking(date time.Time) (ActivePrayerTracking, error) {
	dayPrayers := getDayPrayerTimeFor(date)
	if dayPrayers == nil {
		return ActivePrayerTracking{}, errors.New("Failed to get day prayer")
	}

	previousPrayer, nextPrayer := getNextAndPreviousPrayerTimes(*dayPrayers)
	if previousPrayer == nil || nextPrayer == nil {
		return ActivePrayerTracking{}, errors.New("Could not get previous or next prayer")
	}

	reminaingToNextPrayer := getTimeRemainingTo(nextPrayer.Time)
	if reminaingToNextPrayer == nil {
		return ActivePrayerTracking{}, errors.New("Failed to get time remaining to next prayer")
	}

	timeProgressPercent := timeProgressPercent(previousPrayer.Time, nextPrayer.Time)

	return ActivePrayerTracking{
		DailyPrayerSchedule: DailyPrayerSchedule{
			Date:    dayPrayers.Date,
			Prayers: dayPrayers.Prayers,
		},
		PreviousPrayer: previousPrayer.Name,
		NextPrayer:     nextPrayer.Name,
		TimeRemaining:  *reminaingToNextPrayer,
		Progress:       timeProgressPercent,
	}, nil
}

// getDayPrayerTimeFor get caches data locally or fetch new data from remote then save locally.
// Then search data for specific @year @month and @day. If found return prayer times
func getDayPrayerTimeFor(time time.Time) *DayPrayers {
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
func getNextAndPreviousPrayerTimes(dayPrayers DayPrayers) (*Prayer, *Prayer) {
	day := time.Now()

	yesterdayDate := day.Add(-24 * time.Hour)
	yesterdayPrayers := getDayPrayerTimeFor(yesterdayDate)
	if yesterdayPrayers == nil {
		return nil, nil //errors.New("Failed to get yesterday's prayers")
	}

	tomorrowDate := day.Add(24 * time.Hour)
	tomorrowPrayers := getDayPrayerTimeFor(tomorrowDate)
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

func SameDay(t time.Time, otherT time.Time) bool {
	return t.Year() == otherT.Year() && t.Month() == otherT.Month() && t.Day() == otherT.Day()
}

// getTimeRemainingTo
//
// @Returns
//   - hours remaining
//   - minutes remaining
func getTimeRemainingTo(nextPrayerTime time.Time) *time.Duration {
	now := time.Now()
	if now.After(nextPrayerTime) {
		return nil
	}

	duration := nextPrayerTime.Sub(now)
	return &duration

}

func timeProgressPercent(
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

func formatDate(time time.Time) string {
	day := time.Day()
	month := int(time.Month())
	year := time.Year()

	dayStr := fmt.Sprint(day)
	if day <= 9 {
		dayStr = fmt.Sprintf("0%v", day)
	}
	monthStr := fmt.Sprint(month)
	if month < 9 {
		monthStr = fmt.Sprintf("0%v", month)
	}
	return fmt.Sprintf("%v/%v/%v", dayStr, monthStr, year)
}

func getPrayerTimes(
	data models.PrayerTimesResponse,
	dateStr string,
) *models.DailyPrayersDto {
	for _, dayPrayer := range data.Year {
		if dayPrayer.Gregorian == dateStr {
			return &dayPrayer
		}
	}
	return nil
}
