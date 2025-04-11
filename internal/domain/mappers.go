package domain

import (
	"strings"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

func mapToDayPrayer(prayerTimes models.DailyPrayersDto) *DayPrayers {
	day, err := time.ParseInLocation("02/01/2006", prayerTimes.Gregorian, time.Local)
	if err != nil {
		return nil
	}

	prayers, err := getSortedPrayerTimes(day, prayerTimes.Prayers)
	if err != nil {
		return nil
	}

	return &DayPrayers{
		ID:      prayerTimes.ID,
		Date:    day,
		Prayers: prayers,
	}
}

// getSortedPrayerTimes takes @day
func getSortedPrayerTimes(day time.Time, prayerTimes models.PrayerTimesDto) ([]Prayer, error) {
	result := []Prayer{}

	sortedPrayerNames := models.SortedPrayerNames
	sortedPrayerTimes := prayerTimes.SortedPrayers()
	for i, p := range sortedPrayerTimes {
		t, err := parseTime(day, p)
		if err != nil {
			return []Prayer{}, err
		}
		prayer := Prayer{
			Name: sortedPrayerNames[i],
			Time: t,
		}
		result = append(result, prayer)
	}

	return result, nil
}

// ParseTime takes a time string like this "12:05 pm" and convert it to @time.Time
// or returns error if failed to parse
func parseTime(
	requestedTime time.Time,
	tStr string,
) (time.Time, error) {
	layout := "3:04 pm"
	t, err := time.Parse(layout, strings.ToLower(strings.TrimSpace(tStr)))
	if err != nil {
		return time.Time{}, err
	}
	time := time.Date(requestedTime.Year(), requestedTime.Month(), requestedTime.Day(), t.Hour(), t.Minute(), 0, 0, requestedTime.Location())
	return time, nil
}
