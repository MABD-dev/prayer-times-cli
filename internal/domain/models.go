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

// struct that represents
//   - today's date
//   - today's prayer times
//
// struct that represents
//   - today's date
//   - today's prayer times
//   - time left till next prayer
//   - progress bar from previous to next prayer
