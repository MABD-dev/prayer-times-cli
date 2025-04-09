package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mabd-dev/prayer-times-cli/internal/models"
)

// TODO: Add unit tests to these functions
type Storage interface {
	Save(data models.PrayerTimesResponse) error
	Load(data *models.PrayerTimesResponse) error
}

type FileStorage struct {
	FileName string
}

// Save given data to file
//
// @Returns:
//
//	error if
//	    - Getting file path failed
//	    - marchal data failed
func (s *FileStorage) Save(data models.PrayerTimesResponse) error {
	filePath, err := getOrCreateFilePath((*s).FileName)
	if err != nil {
		return err
	}

	fileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	os.WriteFile(filePath, fileData, 0644)
	return nil
}

// Read/create users json file data dn load them indo data pointer
//
// @Returns:
//
//	error if was not able to get/create file path or read the file
func (s *FileStorage) Load(data *models.PrayerTimesResponse) error {
	filePath, err := getOrCreateFilePath((*s).FileName)
	if err != nil {
		return err
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileData, data)
}

// Create root dir, then a file inside it with given filename
//
// @Returns:
//
//	full path of file: if was able to create it
//	error: if creating root dir failed
func getOrCreateFilePath(filename string) (string, error) {
	prayerTimeCliDir, err := getOrCreateRootDir()
	if err != nil {
		return "", nil
	}
	fileFullPath := filepath.Join(prayerTimeCliDir, filename)
	return fileFullPath, nil
}

// Create hidden root dir for this project at users device home dir
//
// @Returns:
//
//	full path of root dir: if successful
//	error: if getting user home dir failed, or was not able to create folder due to permissions maybe
func getOrCreateRootDir() (string, error) {
	// TODO: make this configurable
	// path from root directory
	prayerTimesCliPath := ".prayer-times-cli"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	prayerTimesCliDir := filepath.Join(homeDir, prayerTimesCliPath)
	if _, err := os.Stat(prayerTimesCliDir); os.IsNotExist(err) {
		err := os.MkdirAll(prayerTimesCliDir, 0700)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return prayerTimesCliDir, nil
}
