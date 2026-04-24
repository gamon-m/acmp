package main

import (
	"acmp/database"
	"acmp/models"
	"acmp/symlink"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	settings models.Settings
	data     database.Data
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	settings, err := getSettings()
	if err != nil {
		panic(err)
	}
	a.settings = settings

	data, err := a.loadData()
	if err != nil {
		panic(err)
	}
	a.data = *data
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {}

func (a *App) GetSettings() models.Settings {
	return a.settings
}

func (a *App) SaveSettings(settings models.Settings) error {
	dbPath, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return err
	}

	err = database.InitSchema(db)
	if err != nil {
		return err
	}

	if settings.ModsPath != a.settings.ModsPath && a.settings.ModsPath != "" {
		err = database.ClearModsAndAutoProfiles(db)
		if err != nil {
			return err
		}
	}

	if settings.AssettoCorsaPath != a.settings.AssettoCorsaPath && a.settings.AssettoCorsaPath != "" {
		existingMods := database.GetModsFromDatabase(db)
		for _, mod := range existingMods {
			linkPath := symlink.BuildSymlinkPath(mod.Category, mod.Name, a.settings.AssettoCorsaPath)
			err := symlink.DeleteSymlink(linkPath)
			if err != nil {
				continue
			}
		}

		err = database.DeactivateAllModsAndProfiles(db)
		if err != nil {
			return err
		}
	}

	err = saveSettings(settings)
	if err != nil {
		return err
	}

	a.settings = settings
	return a.RefreshData()
}

func (a *App) GetData() database.Data {
	return a.data
}

func (a *App) RefreshData() error {
	data, err := a.loadData()
	if err != nil {
		return err
	}
	a.data = *data
	runtime.EventsEmit(a.ctx, "data-updated", nil)
	return nil
}

func (a *App) ResetData() error {
	for _, mod := range a.data.Mods {
		symlinkPath := filepath.Join(a.settings.AssettoCorsaPath, "content", mod.Category, mod.Name)
		err := symlink.DeleteSymlink(symlinkPath)
		if err != nil {
			continue
		}
	}

	a.settings.AssettoCorsaPath = ""
	a.settings.ModsPath = ""
	a.settings.AutomaticProfiles = false

	err := saveSettings(a.settings)
	if err != nil {
		return err
	}

	dbPath, err := getDbPath()
	if err != nil {
		return err
	}
	err = os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	a.data = database.Data{}
	runtime.EventsEmit(a.ctx, "data-updated", nil)
	return nil
}

func (a *App) OpenFolderDialog(folder string) (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: fmt.Sprintf("Select %s folder", folder)})
	return folder, err
}

func getDbPath() (string, error) {
	userDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userDataDir, "acmp", "acmp.db"), nil
}

func (a *App) loadData() (*database.Data, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	err = database.InitSchema(db)
	if err != nil {
		return nil, err
	}

	data := &database.Data{}
	err = data.ScanMods(&a.settings)
	if err != nil {
		return nil, err
	}

	err = data.UpdateDatabase(db, a.settings)
	if err != nil {
		return nil, err
	}

	data.Mods = database.GetModsFromDatabase(db)
	data.Profiles = database.GetProfilesFromDatabase(db)
	data.ModProfiles = database.GetModProfilesFromDatabase(db)

	if a.settings.AssettoCorsaPath != "" {
		err = symlink.ReconcileSymlinks(data.Mods, a.settings.AssettoCorsaPath)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (a *App) SaveProfile(profileRequest models.ProfileRequest) error {
	dbPath, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return err
	}

	var profileId int
	if profileRequest.ProfileId != "" {
		profileId, err = strconv.Atoi(profileRequest.ProfileId)
		if err != nil {
			return err
		}
	} else {
		profileId = -1
	}

	profile := models.Profile{
		Id:          profileId,
		Name:        profileRequest.Name,
		Path:        "",
		Category:    profileRequest.Category,
		Active:      false,
		AutoCreated: false,
	}

	var mods []string
	for _, modDir := range profileRequest.Mods {
		mods = append(mods, modDir)
	}

	if profile.Id == -1 {
		database.CreateProfile(db, profile, mods)
	} else {
		database.UpdateProfile(db, profile, mods)
	}
	return a.RefreshData()
}

func (a *App) DeleteProfile(profileId int) error {
	dbPath, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return err
	}

	err = database.DeleteProfile(db, profileId)
	if err != nil {
		return err
	}
	return a.RefreshData()
}

func (a *App) UpdateProfiles(profiles []models.Profile) error {
	dbPath, err := getDbPath()
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return err
	}

	var profilesToActivate []models.Profile
	for _, profile := range profiles {
		if profile.Active {
			profilesToActivate = append(profilesToActivate, profile)
		}
	}

	var profilesToDeactivate []models.Profile
	for _, profile := range profiles {
		if !profile.Active {
			profilesToDeactivate = append(profilesToDeactivate, profile)
		}
	}

	err = database.ActivateProfiles(db, profilesToActivate)
	if err != nil {
		return err
	}
	err = database.DeactivateProfiles(db, profilesToDeactivate)
	if err != nil {
		return err
	}
	return a.RefreshData()
}
