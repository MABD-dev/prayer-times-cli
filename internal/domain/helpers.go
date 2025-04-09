package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func SameDay(t time.Time, otherT time.Time) bool {
	return t.Year() == otherT.Year() && t.Month() == otherT.Month() && t.Day() == otherT.Day()
}
