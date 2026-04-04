package main

import (
	"acmp/database"
	"acmp/models"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	settings models.Settings
	db       *sql.DB
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

	a.db, err = getDatabase()
	if err != nil {
		panic(err)
	}
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
	err := saveSettings(settings)
	if err != nil {
		return err
	}

	a.settings = settings
	return nil
}

func (a *App) OpenFolderDialog(folder string) (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: fmt.Sprintf("Select %s folder", folder)})
	return folder, err
}

func getAppDataPath() (string, error) {
	userDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userDataDir, "acmp"), nil
}

func getDatabase() (*sql.DB, error) {
	appDataPath, err := getAppDataPath()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDataPath, "acmp.db")
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	err = database.InitSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
