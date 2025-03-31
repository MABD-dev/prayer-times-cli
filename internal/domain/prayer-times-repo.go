package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

// GetDayPrayerTimeFor get caches data locally or fetch new data from remote then save locally.
// Then search data for specific @year @month and @day. If found return prayer times
func GetDayPrayerTimeFor(
	year int,
	month int,
	day int,
) *models.DayPrayerTimes {
	dateStr := formatDate(year, month, day)
	localYearFilename := fmt.Sprintf("%v.json", year)

	data := loadFromLocal(localYearFilename)
	if data == nil {
		fmt.Println("Fetching data from internet...")
		res, err := fetchPrayingTimes(year)
		if err != nil {
			fmt.Println("Failed to fetch data from internet")
			return nil
		}
		Save(localYearFilename, *res)
		data = res
	}

	return getPrayerTimes(*data, dateStr)
}

func loadFromLocal(filename string) *models.PrayerTimesResponse {
	var data models.PrayerTimesResponse
	err := Load(filename, &data)
	if err != nil {
		return nil
	}
	return &data
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

func formatDate(
	year int,
	month int,
	day int,
) string {
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
