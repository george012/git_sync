package config

import (
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_app"
)

var (
	CurrentApp *ExtendApp
)

type ExtendApp struct {
	*gtbox_app.App
	NetListenPortStratumDefault int
	NetListenAPIPortDefault     int
	DualPoolEtcHashInsideUrl    string
}

func NewApp(appName, bundleID, description string, runMode gtbox.RunMode, apiPortDefault int) *ExtendApp {
	app := &ExtendApp{
		App:                         gtbox_app.NewApp(appName, ProjectVersion, bundleID, description, runMode),
		NetListenPortStratumDefault: 0,
		NetListenAPIPortDefault:     apiPortDefault,
	}

	return app
}
