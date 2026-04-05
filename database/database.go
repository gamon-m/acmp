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
		inProfile    BOOLEAN NOT NULL DEFAULT 0,
		lastModified DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS profiles (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT NOT NULL,
		path         TEXT UNIQUE NOT NULL,
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
				if scannedMod.LastModified != dbMod.LastModified ||
					scannedMod.Active != dbMod.Active ||
					scannedMod.InProfile != dbMod.InProfile ||
					scannedMod.Name != dbMod.Name {
					modsToUpdate = append(modsToUpdate, scannedMod)
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
	rows, err := db.Query("SELECT dir, mod_name, category, active, inProfile, lastModified FROM mods")
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

func insertMods(db *sql.DB, mods []models.Mod) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(`INSERT INTO mods (dir, mod_name, category, active, inProfile, lastModified)
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

	statement, err := tx.Prepare(`UPDATE mods SET mod_name = ?, category = ?, active = ?, inProfile = ?, lastModified = ? WHERE dir = ?`)
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
	return tx.Commit()
}
