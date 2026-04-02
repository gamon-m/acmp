package main

import "context"

// App struct
type App struct {
	ctx      context.Context
	settings Settings
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
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {}

func (a *App) GetSettings() Settings {
	return a.settings
}

func (a *App) SaveSettings(settings Settings) error {
	err := saveSettings(settings)
	if err != nil {
		return err
	}

	a.settings = settings
	return nil
}
