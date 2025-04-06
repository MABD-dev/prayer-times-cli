package domain

import (
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

func mapToDayPrayer(prayerTimes models.DayPrayerTimes) *DayPrayers {
	day, err := time.ParseInLocation("02/01/2006", prayerTimes.Gregorian, time.Local)
	if err != nil {
		return nil
	}

	prayers, err := getSortedPrayerTimes(day, prayerTimes.PrayerTimes)
	if err != nil {
		return nil
	}

	return &DayPrayers{
		ID:      prayerTimes.ID,
		Date:    day,
		Prayers: prayers,
	}
}
