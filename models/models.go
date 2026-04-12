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
	Id          int
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

type ProfileRequest struct {
	ProfileId string   `json:"id"`
	Name      string   `json:"name"`
	Category  string   `json:"category"`
	Mods      []string `json:"mods"`
}
