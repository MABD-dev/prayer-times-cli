package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

func loadFromLocal(filename string) *models.PrayerTimesResponse {
	var data models.PrayerTimesResponse
	err := Load(filename, &data)
	if err != nil {
		return nil
	}
	return &data
}

func fetchAndSavePrayerTimes(year int, filename string) (*models.PrayerTimesResponse, error) {
	fmt.Println("Fetching data from internet...")
	fmt.Printf("year=%v\n", year)
	res, err := fetchPrayingTimes(year)
	if err != nil {
		return nil, err
	}
	Save(filename, *res)
	return res, nil
}

func fetchPrayingTimes(year int) (*models.PrayerTimesResponse, error) {
	baseUrl := fmt.Sprintf("https://ibad-al-rahman.github.io/prayer-times/v1/year/days/%v.json", year)
	fmt.Println(baseUrl)

	resp, err := http.Get(baseUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response not ok!!%v\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response models.PrayerTimesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func formatDate(time time.Time) string {
	day := time.Day()
	month := int(time.Month())
	year := time.Year()

	dayStr := fmt.Sprint(day)
	if day < 9 {
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
) *models.DayPrayerTimes {
	for _, dayPrayerTime := range data.Year {
		if dayPrayerTime.Gregorian == dateStr {
			return &dayPrayerTime
		}
	}
	return nil
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
	// Attach today's date to the parsed time
	time := time.Date(requestedTime.Year(), requestedTime.Month(), requestedTime.Day(), t.Hour(), t.Minute(), 59, 0, requestedTime.Location())
	return time, nil
}

// getNextPrayerTime for now assuming their is a next prayer in the day
// TODO: write better docs
func getNextPrayerTime(
	requestedTime time.Time,
	prayerDay time.Time,
	prayerTimes models.PrayerTimes,
) (*time.Time, string) {

	sortedPrayerNames := []string{
		"Fajr",
		"Dhuhu",
		"Asr",
		"Maghrib",
		"Isha",
	}
	sortedPrayerTimes := []string{
		prayerTimes.Fajr,
		prayerTimes.Dhuhr,
		prayerTimes.Asr,
		prayerTimes.Maghrib,
		prayerTimes.Isha,
	}

	for i, prayerTime := range sortedPrayerTimes {
		t, err := parseTime(prayerDay, prayerTime)
		if err != nil {
			return nil, ""
		}

		if t.After(requestedTime) || t.Equal(requestedTime) {
			return &t, sortedPrayerNames[i]
		}
	}

	return nil, ""
}
