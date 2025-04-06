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
