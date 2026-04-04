package main

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func newDatabase(dbPath string) (*sql.DB, error) {
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

func initSchema(db *sql.DB) error {
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
