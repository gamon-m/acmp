package main

import (
	"acmp/models"
	"encoding/json"
	"os"
	"path/filepath"
)

func getSettings() (models.Settings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return models.Settings{}, err
	}

	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = createSettingsFile()
			if err != nil {
				return models.Settings{}, err
			}
			return getDefaultSettings(), nil
		}
		return models.Settings{}, err
	}

	var settings models.Settings
	err = json.Unmarshal(settingsData, &settings)
	if err != nil {
		return models.Settings{}, err
	}

	return settings, nil
}

func saveSettings(settings models.Settings) error {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return err
	}

	settingsData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, settingsData, os.ModePerm)
}

func createSettingsFile() error {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return err
	}

	settingsDir := filepath.Dir(settingsPath)
	err = os.MkdirAll(settingsDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the file with default settings
	settings := getDefaultSettings()
	settingsData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, settingsData, os.ModePerm)
}

func getDefaultSettings() models.Settings {
	return models.Settings{
		AssettoCorsaPath:  "",
		ModsPath:          "",
		AutomaticProfiles: false,
	}
}

func getSettingsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "acmp", "settings.json"), nil
}
