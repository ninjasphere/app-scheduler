package main

import (
	"fmt"

	"github.com/ninjasphere/app-scheduler/controller"
	"github.com/ninjasphere/app-scheduler/model"
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
	scheduler *controller.Scheduler
}

type Config struct {
}

func (a *SchedulerApp) Start(model *model.Schedule) error {
	if a.scheduler != nil {
		return fmt.Errorf("illegal state: scheduler is already running")
	}
	a.scheduler = &controller.Scheduler{}
	err := a.scheduler.Start(model)
	if err == nil {
		a.SendEvent("config", model)
	}
	return err
}

func (a *SchedulerApp) Stop() error {
	var err error
	if a.scheduler != nil {
		tmp := a.scheduler
		a.scheduler = nil
		err = tmp.Stop()
	}
	return err
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
