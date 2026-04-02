package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	AssettoCorsaPath  string `json:"assetto_corsa_path"`
	ModsPath          string `json:"mods_path"`
	AutomaticProfiles bool   `json:"automatic_profiles"`
}

func getSettings() (Settings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return Settings{}, err
	}

	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = createSettingsFile()
			if err != nil {
				return Settings{}, err
			}
			return getDefaultSettings(), nil
		}
		return Settings{}, err
	}

	var settings Settings
	err = json.Unmarshal(settingsData, &settings)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func saveSettings(settings Settings) error {
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

func getDefaultSettings() Settings {
	return Settings{
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
