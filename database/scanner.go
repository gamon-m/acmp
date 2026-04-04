package database

import (
	"acmp/models"
	"io/fs"
	"os"
	"path/filepath"
)

type Data struct {
	Mods        []models.Mod
	Profiles    []models.Profile
	ModProfiles []models.ModProfile
}

func (d *Data) ScanMods(settings *models.Settings) error {
	var category string

	return filepath.WalkDir(settings.ModsPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dir.IsDir() {
			return nil
		}

		if dir.Name() == "tracks" || dir.Name() == "cars" {
			category = dir.Name()
		}

		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".kn5" {
				dirInfo, err := dir.Info()
				if err != nil {
					return err
				}

				mod := models.Mod{
					Dir:          path,
					Name:         filepath.Base(path),
					Category:     category,
					Active:       false,
					InProfile:    false,
					LastModified: dirInfo.ModTime(),
				}
				d.Mods = append(d.Mods, mod)
				return filepath.SkipDir
			}
		}

		return nil
	})
}
