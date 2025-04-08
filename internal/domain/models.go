package domain

import "time"

type Prayer struct {
	Name string
	Time time.Time
}

type DayPrayers struct {
	ID      int
	Date    time.Time
	Prayers []Prayer
}

// To be used in ui

type DailyPrayerSchedule struct {
	Date    time.Time
	Prayers []Prayer
}

type ActivePrayerTracking struct {
	DailyPrayerSchedule
	PreviousPrayer string
	NextPrayer     string
	TimeRemaining  time.Duration
	Progress       float64
}
