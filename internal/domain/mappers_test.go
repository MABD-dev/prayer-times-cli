package domain

import (
	"testing"
	"time"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

func TestMapToDayPrayer(t *testing.T) {
	day := time.Date(2025, 5, 28, 1, 0, 0, 0, time.Now().Location())
	prayers := models.PrayerTimesDto{
		Fajr:    "05:00 am",
		Dhuhr:   "11:00 am",
		Asr:     "01:00 pm",
		Maghrib: "05:00 pm",
		Isha:    "07:00 pm",
	}
	dayPrayers := DayPrayers{
		ID:   1,
		Date: day,
		Prayers: []Prayer{
			{
				Name: "Fajr",
				Time: time.Date(day.Year(), day.Month(), day.Day(), 5, 0, 0, day.Nanosecond(), day.Location()),
			},
			{
				Name: "Dhuhr",
				Time: time.Date(day.Year(), day.Month(), day.Day(), 11, 0, 0, 0, day.Location()),
			},
			{
				Name: "Asr",
				Time: time.Date(day.Year(), day.Month(), day.Day(), 13, 0, 0, 0, day.Location()),
			},
			{
				Name: "Maghrib",
				Time: time.Date(day.Year(), day.Month(), day.Day(), 17, 0, 0, 0, day.Location()),
			},
			{
				Name: "Isha",
				Time: time.Date(day.Year(), day.Month(), day.Day(), 19, 0, 0, 0, day.Location()),
			},
		},
	}

	tests := []struct {
		name               string
		dailyPrayerDto     models.DailyPrayersDto
		expectedDayPrayers *DayPrayers
	}{
		{
			name: "Invalid gregorian date",
			dailyPrayerDto: models.DailyPrayersDto{
				ID:        1,
				Gregorian: "",
				Prayers:   prayers,
			},
			expectedDayPrayers: nil,
		},
		{
			name: "Invalid gregorian date format",
			dailyPrayerDto: models.DailyPrayersDto{
				ID:        1,
				Gregorian: "2025/05/28",
				Prayers:   prayers,
			},
			expectedDayPrayers: nil,
		},
		{
			name: "Invalid gregorian date format 2",
			dailyPrayerDto: models.DailyPrayersDto{
				ID:        1,
				Gregorian: "2025/28/05",
				Prayers:   prayers,
			},
			expectedDayPrayers: nil,
		},
		{
			name: "Valid result",
			dailyPrayerDto: models.DailyPrayersDto{
				ID:        1,
				Gregorian: "28/05/2025",
				Prayers:   prayers,
			},
			expectedDayPrayers: &dayPrayers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToDayPrayer(tt.dailyPrayerDto)

			if tt.expectedDayPrayers == nil && result != nil {
				t.Errorf("Expected dayPrayers to be nil but it was not. result=%v", result)
			}

			if result == nil && tt.expectedDayPrayers == nil {
				return
			}

			if result.ID != tt.expectedDayPrayers.ID {
				t.Errorf("Expected ID=%v, got=%v", tt.expectedDayPrayers.ID, result.ID)
			}

			if !SameDay(result.Date, tt.expectedDayPrayers.Date) {
				t.Errorf("Expected Date=%v, got=%v", tt.expectedDayPrayers.Date, result.Date)
			}

			for i := range 5 {
				if result.Prayers[i].Time.Year() != tt.expectedDayPrayers.Prayers[i].Time.Year() ||
					result.Prayers[i].Time.Month() != tt.expectedDayPrayers.Prayers[i].Time.Month() ||
					result.Prayers[i].Time.Day() != tt.expectedDayPrayers.Prayers[i].Time.Day() ||
					result.Prayers[i].Time.Hour() != tt.expectedDayPrayers.Prayers[i].Time.Hour() ||
					result.Prayers[i].Time.Minute() != tt.expectedDayPrayers.Prayers[i].Time.Minute() ||
					result.Prayers[i].Time.Second() != tt.expectedDayPrayers.Prayers[i].Time.Second() {
					t.Errorf("Time components mismatch\nExpected: Y:%d M:%d D:%d H:%d M:%d S:%d\nGot:      Y:%d M:%d D:%d H:%d M:%d S:%d",
						tt.expectedDayPrayers.Prayers[i].Time.Year(), tt.expectedDayPrayers.Prayers[i].Time.Month(), tt.expectedDayPrayers.Prayers[i].Time.Day(),
						tt.expectedDayPrayers.Prayers[i].Time.Hour(), tt.expectedDayPrayers.Prayers[i].Time.Minute(), tt.expectedDayPrayers.Prayers[i].Time.Second(),
						result.Prayers[i].Time.Year(), result.Prayers[i].Time.Month(), result.Prayers[i].Time.Day(), result.Prayers[i].Time.Hour(), result.Prayers[i].Time.Minute(), result.Prayers[i].Time.Second())
				}

				if result.Prayers[i].Time.Location().String() != tt.expectedDayPrayers.Prayers[i].Time.Location().String() {
					t.Errorf("Location mismatch. Expected %v but got %v", tt.expectedDayPrayers.Prayers[i].Time.Location(), result.Prayers[i].Time.Location())
				}
			}

		})
	}
}

func TestGetSortedPrayerTimes(t *testing.T) {
	day := time.Date(2025, 11, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		day             time.Time
		prayerTimes     models.PrayerTimesDto
		expectedPrayers []Prayer
		expectError     bool
	}{
		{
			name: "Invalid Farj prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr: "",
			},
			expectError: true,
		},
		{
			name: "Invalid Dhuhr prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:  "03:04 pm",
				Dhuhr: "",
			},
			expectError: true,
		},
		{
			name: "Invalid Asr prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:  "03:04 pm",
				Dhuhr: "03:04 pm",
				Asr:   "",
			},
			expectError: true,
		},
		{
			name: "Invalid Maghrib prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:    "03:04 pm",
				Dhuhr:   "03:04 pm",
				Asr:     "03:04 pm",
				Maghrib: "",
			},
			expectError: true,
		},
		{
			name: "Invalid Isha prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:    "03:04 pm",
				Dhuhr:   "03:04 pm",
				Asr:     "03:04 pm",
				Maghrib: "03:04 pm",
				Isha:    "",
			},
			expectError: true,
		},
		{
			name: "Invalid Isha prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:    "14:04 pm",
				Dhuhr:   "03:04 pm",
				Asr:     "03:04 pm",
				Maghrib: "03:04 pm",
				Isha:    "03:04 pm",
			},
			expectError: true,
		},
		{
			name: "Valid prayer time",
			day:  day,
			prayerTimes: models.PrayerTimesDto{
				Fajr:    "05:00 am",
				Dhuhr:   "11:00 am",
				Asr:     "01:00 pm",
				Maghrib: "05:00 pm",
				Isha:    "07:00 pm",
			},
			expectedPrayers: []Prayer{
				{
					Name: "Fajr",
					Time: time.Date(day.Year(), day.Month(), day.Day(), 5, 0, 0, day.Nanosecond(), day.Location()),
				},
				{
					Name: "Dhuhr",
					Time: time.Date(day.Year(), day.Month(), day.Day(), 11, 0, 0, 0, day.Location()),
				},
				{
					Name: "Asr",
					Time: time.Date(day.Year(), day.Month(), day.Day(), 13, 0, 0, 0, day.Location()),
				},
				{
					Name: "Maghrib",
					Time: time.Date(day.Year(), day.Month(), day.Day(), 17, 0, 0, 0, day.Location()),
				},
				{
					Name: "Isha",
					Time: time.Date(day.Year(), day.Month(), day.Day(), 19, 0, 0, 0, day.Location()),
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getSortedPrayerTimes(tt.day, tt.prayerTimes)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if !tt.expectError {
				if len(result) != 5 {
					t.Errorf("Expected 5 prayers got: %v", len(result))
				}

				for i := range 5 {
					if result[i].Time.Year() != tt.expectedPrayers[i].Time.Year() ||
						result[i].Time.Month() != tt.expectedPrayers[i].Time.Month() ||
						result[i].Time.Day() != tt.expectedPrayers[i].Time.Day() ||
						result[i].Time.Hour() != tt.expectedPrayers[i].Time.Hour() ||
						result[i].Time.Minute() != tt.expectedPrayers[i].Time.Minute() ||
						result[i].Time.Second() != tt.expectedPrayers[i].Time.Second() {
						t.Errorf("Time components mismatch\nExpected: Y:%d M:%d D:%d H:%d M:%d S:%d\nGot:      Y:%d M:%d D:%d H:%d M:%d S:%d",
							tt.expectedPrayers[i].Time.Year(), tt.expectedPrayers[i].Time.Month(), tt.expectedPrayers[i].Time.Day(),
							tt.expectedPrayers[i].Time.Hour(), tt.expectedPrayers[i].Time.Minute(), tt.expectedPrayers[i].Time.Second(),
							result[i].Time.Year(), result[i].Time.Month(), result[i].Time.Day(), result[i].Time.Hour(), result[i].Time.Minute(), result[i].Time.Second())
					}

					if result[i].Time.Location().String() != tt.expectedPrayers[i].Time.Location().String() {
						t.Errorf("Location mismatch. Expected %v but got %v", tt.expectedPrayers[i].Time.Location(), result[i].Time.Location())
					}
				}
			}
		})
	}

}

func TestParseTime(t *testing.T) {
	// Setup a known base time for consistency in tests
	baseTime := time.Date(2023, 11, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		requestedTime time.Time
		timeStr       string
		expected      time.Time
		expectError   bool
	}{
		{
			name:          "Valid morning time",
			requestedTime: baseTime,
			timeStr:       "5:30 am",
			expected:      time.Date(2023, 11, 15, 5, 30, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Valid afternoon time",
			requestedTime: baseTime,
			timeStr:       "2:45 pm",
			expected:      time.Date(2023, 11, 15, 14, 45, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with extra whitespace",
			requestedTime: baseTime,
			timeStr:       "  7:15 pm  ",
			expected:      time.Date(2023, 11, 15, 19, 15, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with uppercase meridian",
			requestedTime: baseTime,
			timeStr:       "9:00 AM",
			expected:      time.Date(2023, 11, 15, 9, 0, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with mixed case meridian",
			requestedTime: baseTime,
			timeStr:       "11:30 Pm",
			expected:      time.Date(2023, 11, 15, 23, 30, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Midnight",
			requestedTime: baseTime,
			timeStr:       "12:00 am",
			expected:      time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Noon",
			requestedTime: baseTime,
			timeStr:       "12:00 pm",
			expected:      time.Date(2023, 11, 15, 12, 0, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Different day in requested time",
			requestedTime: time.Date(2024, 2, 29, 15, 0, 0, 0, time.UTC), // Leap year
			timeStr:       "8:45 am",
			expected:      time.Date(2024, 2, 29, 8, 45, 0, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Different time zone",
			requestedTime: time.Date(2023, 11, 15, 10, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			timeStr:       "3:30 pm",
			expected:      time.Date(2023, 11, 15, 15, 30, 0, 0, time.FixedZone("EST", -5*60*60)),
			expectError:   false,
		},
		{
			name:          "Invalid time format",
			requestedTime: baseTime,
			timeStr:       "25:70 pm",
			expectError:   true,
		},
		{
			name:          "Invalid meridian",
			requestedTime: baseTime,
			timeStr:       "10:30 xx",
			expectError:   true,
		},
		{
			name:          "Empty string",
			requestedTime: baseTime,
			timeStr:       "",
			expectError:   true,
		},
		{
			name:          "Invalid format - missing meridian",
			requestedTime: baseTime,
			timeStr:       "10:30",
			expectError:   true,
		},
		{
			name:          "Invalid format - extra text",
			requestedTime: baseTime,
			timeStr:       "10:30 pm extra",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTime(tt.requestedTime, tt.timeStr)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError {
				if !result.Equal(tt.expected) {
					t.Errorf("Expected %v but got %v", tt.expected, result)
				}

				if result.Year() != tt.expected.Year() ||
					result.Month() != tt.expected.Month() ||
					result.Day() != tt.expected.Day() ||
					result.Hour() != tt.expected.Hour() ||
					result.Minute() != tt.expected.Minute() ||
					result.Second() != tt.expected.Second() {
					t.Errorf("Time components mismatch\nExpected: Y:%d M:%d D:%d H:%d M:%d S:%d\nGot:      Y:%d M:%d D:%d H:%d M:%d S:%d",
						tt.expected.Year(), tt.expected.Month(), tt.expected.Day(), tt.expected.Hour(), tt.expected.Minute(), tt.expected.Second(),
						result.Year(), result.Month(), result.Day(), result.Hour(), result.Minute(), result.Second())
				}

				if result.Location().String() != tt.expected.Location().String() {
					t.Errorf("Location mismatch. Expected %v but got %v", tt.expected.Location(), result.Location())
				}
			}
		})
	}
}
