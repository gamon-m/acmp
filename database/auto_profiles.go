package database

import (
	"acmp/models"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

func syncAutoProfiles(db *sql.DB, mods []models.Mod, modsPath string) error {
	existingProfiles := GetProfilesFromDatabase(db)
	existingModProfiles := GetModProfilesFromDatabase(db)

	expectedProfiles := deriveExpectedProfiles(mods, modsPath)

	var profilesToInsert []models.Profile
	var profilesToDelete []int

	existingAutoProfilesByPath := make(map[string]models.Profile)
	for _, profile := range existingProfiles {
		if profile.AutoCreated {
			existingAutoProfilesByPath[profile.Path] = profile
		}
	}

	for path, profile := range expectedProfiles {
		if _, exists := existingAutoProfilesByPath[path]; !exists {
			profilesToInsert = append(profilesToInsert, profile)
		}
	}

	for path, existing := range existingAutoProfilesByPath {
		if _, expected := expectedProfiles[path]; !expected {
			profilesToDelete = append(profilesToDelete, existing.Id)
		}
	}

	insertProfiles(db, profilesToInsert)

	allProfiles := GetProfilesFromDatabase(db)
	profileMapByPath := make(map[string]models.Profile)
	for _, profile := range allProfiles {
		profileMapByPath[profile.Path] = profile
	}

	expectedModProfiles := deriveExpectedModProfiles(mods, profileMapByPath, modsPath)

	existingModProfilesSet := make(map[string]bool)
	for _, modProfile := range existingModProfiles {
		key := fmt.Sprintf("%s:%d", modProfile.ModDir, modProfile.ProfileId)
		existingModProfilesSet[key] = true
	}

	var modProfilesToInsert []models.ModProfile
	for _, expected := range expectedModProfiles {
		key := fmt.Sprintf("%s:%d", expected.ModDir, expected.ProfileId)
		if !existingModProfilesSet[key] {
			modProfilesToInsert = append(modProfilesToInsert, expected)
		}
	}

	insertModProfiles(db, modProfilesToInsert)
	deleteProfiles(db, profilesToDelete)

	return nil
}

func deriveExpectedProfiles(mods []models.Mod, modsPath string) map[string]models.Profile {
	profileMap := make(map[string]models.Profile)

	for _, mod := range mods {
		relPath, err := filepath.Rel(modsPath, mod.Dir)
		if err != nil {
			continue
		}
		relParts := strings.Split(filepath.ToSlash(relPath), "/")
		for i := 1; i < len(relParts)-1; i++ {
			profilePath := strings.Join(relParts[:i+1], "/")
			profileName := strings.Join(relParts[1:i+1], "/")

			if _, exists := profileMap[profilePath]; exists {
				continue
			}

			profileMap[profilePath] = models.Profile{
				Name:        profileName,
				Path:        profilePath,
				Category:    mod.Category,
				Active:      false,
				AutoCreated: true,
			}
		}
	}
	return profileMap
}

func deriveExpectedModProfiles(mods []models.Mod, profileMap map[string]models.Profile, modsPath string) []models.ModProfile {
	var expected []models.ModProfile

	for _, mod := range mods {
		relPath, err := filepath.Rel(modsPath, mod.Dir)
		if err != nil {
			continue
		}
		relParts := strings.Split(filepath.ToSlash(relPath), "/")

		for i := 1; i < len(relParts)-1; i++ {
			profilePath := strings.Join(relParts[:i+1], "/")
			if profile, exists := profileMap[profilePath]; exists {
				expected = append(expected, models.ModProfile{
					ModDir:    mod.Dir,
					ProfileId: profile.Id,
				})
			}
		}
	}
	return expected
}

func insertProfiles(db *sql.DB, profiles []models.Profile) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`INSERT INTO profiles (name, path, category, active, auto_created) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, profile := range profiles {
		_, err := statement.Exec(profile.Name, profile.Path, profile.Category, profile.Active, profile.AutoCreated)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func insertModProfiles(db *sql.DB, modProfiles []models.ModProfile) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`INSERT INTO mod_profiles (mod_dir, profile_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, modProfile := range modProfiles {
		_, err := statement.Exec(modProfile.ModDir, modProfile.ProfileId)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func deleteProfiles(db *sql.DB, profileIds []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteProfileMods, err := tx.Prepare(`DELETE FROM mod_profiles WHERE profile_id = ?`)
	if err != nil {
		return err
	}
	defer deleteProfileMods.Close()
	for _, id := range profileIds {
		_, err := deleteProfileMods.Exec(id)
		if err != nil {
			return err
		}
	}

	deleteProfiles, err := tx.Prepare(`DELETE FROM profiles WHERE id = ?`)
	if err != nil {
		return err
	}
	defer deleteProfiles.Close()
	for _, id := range profileIds {
		_, err := deleteProfiles.Exec(id)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
