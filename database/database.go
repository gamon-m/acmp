package database

import (
	"acmp/models"
	"acmp/symlink"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func NewDatabase(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS mods (
		dir          TEXT PRIMARY KEY,
		mod_name     TEXT NOT NULL,
		category     TEXT NOT NULL,
		active       BOOLEAN NOT NULL DEFAULT 0,
		in_profile    BOOLEAN NOT NULL DEFAULT 0,
		last_modified DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS profiles (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT NOT NULL,
		path         TEXT,
		category     TEXT NOT NULL,
		active       BOOLEAN NOT NULL DEFAULT 0,
		auto_created BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS mod_profiles (
		mod_dir    TEXT NOT NULL,
		profile_id INTEGER NOT NULL,
		PRIMARY KEY (mod_dir, profile_id),
		FOREIGN KEY (mod_dir) REFERENCES mods(dir),
		FOREIGN KEY (profile_id) REFERENCES profiles(id)
	);

	CREATE INDEX IF NOT EXISTS idx_mod_profiles_mod ON mod_profiles(mod_dir);
	CREATE INDEX IF NOT EXISTS idx_mod_profiles_profile ON mod_profiles(profile_id);
	`

	_, err := db.Exec(schema)
	return err
}

func (d *Data) UpdateDatabase(db *sql.DB, settings models.Settings) error {
	if settings.AutomaticProfiles {
		err := syncAutoProfiles(db, d.Mods, settings.ModsPath)
		if err != nil {
			return err
		}
	} else {
		err := deleteAutoProfiles(db)
		if err != nil {
			return err
		}
	}

	databaseMods := GetModsFromDatabase(db)
	modProfiles := GetModProfilesFromDatabase(db)

	var modsToAdd []models.Mod
	var modsToUpdate []models.Mod
	var modsToDelete []models.Mod

	modsInProfilesSet := make(map[string]bool)
	for _, mp := range modProfiles {
		modsInProfilesSet[mp.ModDir] = true
	}

	for _, scannedMod := range d.Mods {
		dbMod, found := findModByDir(databaseMods, scannedMod.Dir)
		if !found {
			modsToAdd = append(modsToAdd, scannedMod)
			continue
		}

		inProfile := modsInProfilesSet[scannedMod.Dir]
		needsUpdate := scannedMod.LastModified != dbMod.LastModified || inProfile != dbMod.InProfile

		if needsUpdate {
			modsToUpdate = append(modsToUpdate, models.Mod{
				Dir:          scannedMod.Dir,
				Name:         scannedMod.Name,
				Category:     scannedMod.Category,
				Active:       dbMod.Active,
				InProfile:    inProfile,
				LastModified: scannedMod.LastModified,
			})
		}
	}

	for _, dbMod := range databaseMods {
		if !modExists(d.Mods, dbMod.Dir) {
			modsToDelete = append(modsToDelete, dbMod)
		}
	}

	insertMods(db, modsToAdd)
	updateMods(db, modsToUpdate)
	deleteMods(db, modsToDelete, settings.AssettoCorsaPath)
	return nil
}

func findModByDir(mods []models.Mod, dir string) (models.Mod, bool) {
	for _, mod := range mods {
		if mod.Dir == dir {
			return mod, true
		}
	}
	return models.Mod{}, false
}

func modExists(mods []models.Mod, dir string) bool {
	for _, mod := range mods {
		if mod.Dir == dir {
			return true
		}
	}
	return false
}

func GetModsFromDatabase(db *sql.DB) []models.Mod {
	rows, err := db.Query("SELECT dir, mod_name, category, active, in_profile, last_modified FROM mods")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var mods []models.Mod
	for rows.Next() {
		var mod models.Mod
		err := rows.Scan(&mod.Dir, &mod.Name, &mod.Category, &mod.Active, &mod.InProfile, &mod.LastModified)
		if err != nil {
			continue
		}
		mods = append(mods, mod)
	}
	return mods
}

func GetProfilesFromDatabase(db *sql.DB) []models.Profile {
	rows, err := db.Query("SELECT id, name, path, category, active, auto_created FROM profiles")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile
		err := rows.Scan(&profile.Id, &profile.Name, &profile.Path, &profile.Category, &profile.Active, &profile.AutoCreated)
		if err != nil {
			continue
		}
		profiles = append(profiles, profile)
	}
	return profiles
}

func GetModProfilesFromDatabase(db *sql.DB) []models.ModProfile {
	rows, err := db.Query("SELECT mod_dir, profile_id FROM mod_profiles")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var modProfiles []models.ModProfile
	for rows.Next() {
		var modProfile models.ModProfile
		err := rows.Scan(&modProfile.ModDir, &modProfile.ProfileId)
		if err != nil {
			continue
		}
		modProfiles = append(modProfiles, modProfile)
	}
	return modProfiles
}

func CreateProfile(db *sql.DB, profile models.Profile, modDirs []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO profiles (name, path, category, active, auto_created) VALUES (?, ?, ?, ?, ?)`,
		profile.Name, profile.Path, profile.Category, profile.Active, profile.AutoCreated)
	if err != nil {
		return err
	}

	profileId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	statement, err := tx.Prepare(`INSERT INTO mod_profiles (mod_dir, profile_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, modDir := range modDirs {
		_, err := statement.Exec(modDir, profileId)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func UpdateProfile(db *sql.DB, profile models.Profile, modDirs []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE profiles SET name = ? WHERE id = ?`, profile.Name, profile.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM mod_profiles where profile_id = ?`, profile.Id)
	if err != nil {
		return err
	}

	statement, err := tx.Prepare(`INSERT INTO mod_profiles (mod_dir, profile_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, modDir := range modDirs {
		_, err := statement.Exec(modDir, profile.Id)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func DeleteProfile(db *sql.DB, profileId int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM mod_profiles WHERE profile_id = ?`, profileId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM profiles WHERE id = ?`, profileId)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func ActivateProfiles(db *sql.DB, profiles []models.Profile) error {
	var mods []string

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, profile := range profiles {
		modDirs := getModsInProfile(db, profile.Id)
		mods = append(mods, modDirs...)
		_, err := tx.Exec(`UPDATE profiles SET active = 1 WHERE id = ?`, profile.Id)
		if err != nil {
			return err
		}
	}

	for _, modDir := range mods {
		_, err := tx.Exec(`UPDATE mods SET active = 1 WHERE dir = ?`, modDir)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func DeactivateProfiles(db *sql.DB, profiles []models.Profile) error {
	var mods []string

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, profile := range profiles {
		modDirs := getModsInProfile(db, profile.Id)
		mods = append(mods, modDirs...)
		_, err := tx.Exec(`UPDATE profiles SET active = 0 WHERE id = ?`, profile.Id)
		if err != nil {
			return err
		}
	}

	filteredMods := make(map[string]struct{})
	for _, modDir := range mods {
		if modInOtherActiveProfiles(db, modDir, profiles) {
			continue
		}
		filteredMods[modDir] = struct{}{}
	}

	for modDir := range filteredMods {
		_, err := tx.Exec(`UPDATE mods SET active = 0 WHERE dir = ?`, modDir)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func modInOtherActiveProfiles(db *sql.DB, modDir string, profiles []models.Profile) bool {
	var ids []int
	for _, profile := range profiles {
		ids = append(ids, profile.Id)
	}

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(
		"SELECT id FROM profiles WHERE id NOT IN (%s) AND active = 1 AND id IN (SELECT profile_id FROM mod_profiles WHERE mod_dir = ?)",
		strings.Join(placeholders, ","),
	)

	args := make([]interface{}, len(ids)+1)
	for i, id := range ids {
		args[i] = id
	}
	args[len(ids)] = modDir

	rows, err := db.Query(query, args...)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func getModsInProfile(db *sql.DB, profileId int) []string {
	rows, err := db.Query(`SELECT mod_dir FROM mod_profiles WHERE profile_id = ?`, profileId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var modDirs []string
	for rows.Next() {
		var modDir string
		err := rows.Scan(&modDir)
		if err != nil {
			continue
		}
		modDirs = append(modDirs, modDir)
	}
	return modDirs
}

func insertMods(db *sql.DB, mods []models.Mod) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`INSERT INTO mods (dir, mod_name, category, active, in_profile, last_modified)
	 VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, mod := range mods {
		_, err := statement.Exec(mod.Dir, mod.Name, mod.Category, mod.Active, mod.InProfile, mod.LastModified)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func updateMods(db *sql.DB, mods []models.Mod) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`UPDATE mods SET mod_name = ?, category = ?, active = ?, in_profile = ?, last_modified = ? WHERE dir = ?`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, mod := range mods {
		_, err := statement.Exec(mod.Name, mod.Category, mod.Active, mod.InProfile, mod.LastModified, mod.Dir)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func deleteMods(db *sql.DB, mods []models.Mod, assettoCorsaPath string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`DELETE FROM mods WHERE dir = ?`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, mod := range mods {
		linkPath := symlink.BuildSymlinkPath(mod.Category, mod.Name, assettoCorsaPath)
		err := symlink.DeleteSymlink(linkPath)
		if err != nil {
			log.Printf("warning: failed to delete symlink for %s: %v", mod.Name, err)
		}

		_, err = statement.Exec(mod.Dir)
		if err != nil {
			return err
		}
	}

	statement, err = tx.Prepare(`DELETE FROM mod_profiles WHERE mod_dir = ?`)
	if err != nil {
		return err
	}
	defer statement.Close()
	for _, mod := range mods {
		_, err := statement.Exec(mod.Dir)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func deleteAutoProfiles(db *sql.DB) error {
	profiles := GetProfilesFromDatabase(db)
	for _, profile := range profiles {
		if profile.AutoCreated {
			err := DeleteProfile(db, profile.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ClearModsAndAutoProfiles(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM mod_profiles`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM mods`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM profiles WHERE auto_created = 1`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE profiles SET active = 0`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeactivateAllModsAndProfiles(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE mods SET active = 0`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE profiles SET active = 0`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
