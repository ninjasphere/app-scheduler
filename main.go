package main

import (
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

type service struct {
}

type SchedulerApp struct {
	support.AppSupport
}

type Config struct {
}

func (a *SchedulerApp) Start(config *Config) error {
	a.SendEvent("config", config)
	return nil
}

func (a *SchedulerApp) Stop() error {
	return nil
}

func main() {
	app := &SchedulerApp{}
	err := app.Init(info)
	if err != nil {
		app.Log.Fatalf("failed to initialize app: %v", err)
	}

	err = app.Export(app)
	if err != nil {
		app.Log.Fatalf("failed to export app: %v", err)
	}

	support.WaitUntilSignal()
}
