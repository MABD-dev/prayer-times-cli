package models

type Prayers struct {
	Fajr    string `json:"fajr"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"ishaa"`
}

type Event struct {
	En string `json:"en"`
	Ar string `json:"ar"`
}

type DayPrayers struct {
	ID        int     `json:"id"`
	WeekID    int     `json:"weekId"`
	Gregorian string  `json:"gregorian"`
	Hijri     string  `json:"hijri"`
	Prayers   Prayers `json:"prayerTimes"`
	Event     Event   `json:"event"`
}

type PrayerTimesResponse struct {
	Year []DayPrayers `json:"year"`
	Sha1 string       `json:"sha1"`
}

// SortedPrayers return list of prayer times in ascending order
func (p Prayers) SortedPrayers() []string {
	return []string{
		p.Fajr,
		p.Dhuhr,
		p.Asr,
		p.Maghrib,
		p.Isha,
	}
}

// This should be the only place where we define prayer names as
// strings. To be used through out the app
var (
	SortedPrayerNames = []string{
		"Fajr",
		"Dhuhr",
		"Asr",
		"Maghrib",
		"Isha",
	}
)
