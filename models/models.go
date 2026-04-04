package models

import (
	"time"
)

type Mod struct {
	Dir          string
	Name         string
	Category     string
	Active       bool
	InProfile    bool
	LastModified time.Time
}

type Profile struct {
	ID          int
	Name        string
	Path        string
	Category    string
	Active      bool
	AutoCreated bool
}

type ModProfile struct {
	ModDir    string
	ProfileID int
}

type Settings struct {
	AssettoCorsaPath  string `json:"assetto_corsa_path"`
	ModsPath          string `json:"mods_path"`
	AutomaticProfiles bool   `json:"automatic_profiles"`
}
