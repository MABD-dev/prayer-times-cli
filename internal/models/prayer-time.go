package models

type PrayerTimes struct {
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

type DayPrayerTimes struct {
	ID          int         `json:"id"`
	WeekID      int         `json:"weekId"`
	Gregorian   string      `json:"gregorian"`
	Hijri       string      `json:"hijri"`
	PrayerTimes PrayerTimes `json:"prayerTimes"`
	Event       Event       `json:"event"`
}

type PrayerTimesResponse struct {
	Year []DayPrayerTimes `json:"year"`
	Sha1 string           `json:"sha1"`
}
