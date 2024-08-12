package config

import (
	"fmt"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_app"
)

var (
	CurrentApp *ExtendApp
)

type ExtendApp struct {
	*gtbox_app.App
	APIListenPort int
	ReposDir      string
}

func NewApp(appName, bundleID, description string, runMode gtbox.RunMode, apiPortDefault int) *ExtendApp {
	app := &ExtendApp{
		App:           gtbox_app.NewApp(appName, ProjectVersion, bundleID, description, runMode),
		APIListenPort: apiPortDefault,
	}

	app.ReposDir = fmt.Sprintf("%s/repos_handler", app.AppRunAsDir)
	return app
}
