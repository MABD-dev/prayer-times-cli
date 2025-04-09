package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetOrCreateRootDir tests the getOrCreateRootDir function
func TestGetOrCreateRootDir(t *testing.T) {
	// Test successful directory creation
	dir, err := getOrCreateRootDir()
	require.NoError(t, err, "Expected no error from getOrCreateRootDir")
	require.NotEmpty(t, dir, "Expected non-empty directory path")

	// Verify directory exists
	info, err := os.Stat(dir)
	require.NoError(t, err, "Directory should exist")
	require.True(t, info.IsDir(), "Path should be a directory")

	// Test permissions
	require.Equal(t, os.FileMode(0700), info.Mode().Perm()&0700, "Directory should have 0700 permissions")
}

// TestGetOrCreateFilePath tests the getOrCreateFilePath function
func TestGetOrCreateFilePath(t *testing.T) {
	// Test with valid filename
	testFilename := "test-file.json"
	filePath, err := getOrCreateFilePath(testFilename)
	require.NoError(t, err, "Expected no error from getOrCreateFilePath")
	require.NotEmpty(t, filePath, "Expected non-empty file path")

	// Verify the path ends with the correct filename
	require.Equal(t, testFilename, filepath.Base(filePath), "File path should end with the specified filename")
}

// TestFileStorageSave tests the Save method of FileStorage
func TestFileStorageSave(t *testing.T) {
	// Create test data matching the PrayerTimesResponse structure
	testData := models.PrayerTimesResponse{
		Year: []models.DailyPrayersDto{
			{
				ID:        1,
				WeekID:    1,
				Gregorian: "2024-04-09",
				Hijri:     "1445-09-30",
				Prayers: models.PrayerTimesDto{
					Fajr:    "05:30",
					Dhuhr:   "12:45",
					Asr:     "16:15",
					Maghrib: "19:30",
					Isha:    "21:00",
				},
				Event: models.Event{
					En: "Test Event",
					Ar: "حدث اختبار",
				},
			},
		},
		Sha1: "test-sha1-hash",
	}

	// Create a temporary file for testing
	tempFile := "test-save.json"
	storage := FileStorage{FileName: tempFile}

	// Test saving
	err := storage.Save(testData)
	require.NoError(t, err, "Save should not return an error")

	// Verify the file exists
	filePath, err := getOrCreateFilePath(tempFile)
	require.NoError(t, err, "Getting file path should not fail")

	fileInfo, err := os.Stat(filePath)
	require.NoError(t, err, "File should exist")
	require.False(t, fileInfo.IsDir(), "Should be a file, not a directory")

	// Verify file contents
	fileData, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file")

	var savedData models.PrayerTimesResponse
	err = json.Unmarshal(fileData, &savedData)
	require.NoError(t, err, "File should contain valid JSON")
	assert.Equal(t, testData, savedData, "Saved data should match the original")

	// Cleanup
	os.Remove(filePath)
}

// TestFileStorageLoad tests the Load method of FileStorage
func TestFileStorageLoad(t *testing.T) {
	// Create test data matching the PrayerTimesResponse structure
	testData := models.PrayerTimesResponse{
		Year: []models.DailyPrayersDto{
			{
				ID:        1,
				WeekID:    1,
				Gregorian: "2024-04-09",
				Hijri:     "1445-09-30",
				Prayers: models.PrayerTimesDto{
					Fajr:    "05:30",
					Dhuhr:   "12:45",
					Asr:     "16:15",
					Maghrib: "19:30",
					Isha:    "21:00",
				},
				Event: models.Event{
					En: "Test Event",
					Ar: "حدث اختبار",
				},
			},
		},
		Sha1: "test-sha1-hash",
	}

	// Create a temporary file with test data
	tempFile := "test-load.json"
	storage := FileStorage{FileName: tempFile}

	// Save the test data first
	err := storage.Save(testData)
	require.NoError(t, err, "Save should not return an error")

	// Test loading
	var loadedData models.PrayerTimesResponse
	err = storage.Load(&loadedData)
	require.NoError(t, err, "Load should not return an error")
	assert.Equal(t, testData, loadedData, "Loaded data should match the original")

	// Cleanup
	filePath, _ := getOrCreateFilePath(tempFile)
	os.Remove(filePath)
}

// TestFileStorageLoadNonExistent tests the Load method with a non-existent file
func TestFileStorageLoadNonExistent(t *testing.T) {
	// Use a filename that doesn't exist
	storage := FileStorage{FileName: "non-existent-file.json"}

	var loadedData models.PrayerTimesResponse
	err := storage.Load(&loadedData)
	require.Error(t, err, "Load should return an error for non-existent file")
}

// TestGetOrCreateFilePath_Error tests error handling in getOrCreateFilePath
func TestGetOrCreateFilePath_Error(t *testing.T) {
	// For a proper test, we would need to mock the getOrCreateRootDir function
	// to return an error, which would require dependency injection or
	// function variable techniques
	t.Skip("Skipping test for getOrCreateFilePath error - would require mocking")
}

// TestGetOrCreateRootDir_UserHomeDirError tests error handling when os.UserHomeDir fails
func TestGetOrCreateRootDir_UserHomeDirError(t *testing.T) {
	// For a proper test, we would need to mock the os.UserHomeDir function
	// to return an error, which would require dependency injection or
	// function variable techniques
	t.Skip("Skipping test for os.UserHomeDir error - would require mocking the function")
}

// TestGetOrCreateRootDir_MkdirError tests error handling when os.MkdirAll fails
func TestGetOrCreateRootDir_MkdirError(t *testing.T) {
	// For a proper test, we would need to mock the os.MkdirAll function
	// to return an error, which would require dependency injection or
	// function variable techniques
	t.Skip("Skipping test for os.MkdirAll error - would require mocking the function")
}
