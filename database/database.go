package database

import (
	"acmp/models"
	"database/sql"
	"os"
	"path/filepath"

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

func (d *Data) UpdateDatabase(db *sql.DB) error {
	databaseMods := GetModsFromDatabase(db)
	var modsToAdd []models.Mod
	var modsToUpdate []models.Mod
	var modsToDelete []models.Mod

	for _, scannedMod := range d.Mods {
		found := false
		for _, dbMod := range databaseMods {
			if scannedMod.Dir == dbMod.Dir {
				found = true
				if scannedMod.LastModified != dbMod.LastModified {
					modsToUpdate = append(modsToUpdate, models.Mod{
						Dir:          scannedMod.Dir,
						Name:         scannedMod.Name,
						Category:     scannedMod.Category,
						Active:       dbMod.Active,
						InProfile:    dbMod.InProfile,
						LastModified: scannedMod.LastModified,
					})
				}
				break
			}
		}
		if !found {
			modsToAdd = append(modsToAdd, scannedMod)
		}
	}

	for _, dbMod := range databaseMods {
		found := false
		for _, scannedMod := range d.Mods {
			if dbMod.Dir == scannedMod.Dir {
				found = true
				break
			}
		}
		if !found {
			modsToDelete = append(modsToDelete, dbMod)
		}
	}

	insertMods(db, modsToAdd)
	updateMods(db, modsToUpdate)
	deleteMods(db, modsToDelete)

	return nil
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
		for _, modDir := range modDirs {
			if modInOtherActiveProfiles(db, modDir, profile.Id) {
				continue
			}
			mods = append(mods, modDir)
		}
		_, err := tx.Exec(`UPDATE profiles SET active = 0 WHERE id = ?`, profile.Id)
		if err != nil {
			return err
		}
	}

	for _, modDir := range mods {
		_, err := tx.Exec(`UPDATE mods SET active = 0 WHERE dir = ?`, modDir)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func modInOtherActiveProfiles(db *sql.DB, modDir string, profileId int) bool {
	rows, err := db.Query(`SELECT id FROM profiles WHERE id != ? AND active = 1 AND id IN (SELECT profile_id FROM mod_profiles WHERE mod_dir = ?)`, profileId, modDir)
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

func deleteMods(db *sql.DB, mods []models.Mod) error {
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
		_, err := statement.Exec(mod.Dir)
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
