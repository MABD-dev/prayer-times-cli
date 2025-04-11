package domain

import (
	"testing"
	"time"
)

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
			expected:      time.Date(2023, 11, 15, 5, 30, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Valid afternoon time",
			requestedTime: baseTime,
			timeStr:       "2:45 pm",
			expected:      time.Date(2023, 11, 15, 14, 45, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with extra whitespace",
			requestedTime: baseTime,
			timeStr:       "  7:15 pm  ",
			expected:      time.Date(2023, 11, 15, 19, 15, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with uppercase meridian",
			requestedTime: baseTime,
			timeStr:       "9:00 AM",
			expected:      time.Date(2023, 11, 15, 9, 0, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Time with mixed case meridian",
			requestedTime: baseTime,
			timeStr:       "11:30 Pm",
			expected:      time.Date(2023, 11, 15, 23, 30, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Midnight",
			requestedTime: baseTime,
			timeStr:       "12:00 am",
			expected:      time.Date(2023, 11, 15, 0, 0, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Noon",
			requestedTime: baseTime,
			timeStr:       "12:00 pm",
			expected:      time.Date(2023, 11, 15, 12, 0, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Different day in requested time",
			requestedTime: time.Date(2024, 2, 29, 15, 0, 0, 0, time.UTC), // Leap year
			timeStr:       "8:45 am",
			expected:      time.Date(2024, 2, 29, 8, 45, 59, 0, time.UTC),
			expectError:   false,
		},
		{
			name:          "Different time zone",
			requestedTime: time.Date(2023, 11, 15, 10, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			timeStr:       "3:30 pm",
			expected:      time.Date(2023, 11, 15, 15, 30, 59, 0, time.FixedZone("EST", -5*60*60)),
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
