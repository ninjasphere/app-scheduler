package main

import (
	"fmt"

	"github.com/ninjasphere/app-scheduler/controller"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/rpc"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type SchedulerApp struct {
	support.AppSupport
	scheduler *controller.Scheduler
	model     *model.Schedule
	service   *rpc.ExportedService
}

// Start is called after the ExportApp call is complete.
func (a *SchedulerApp) Start(model *model.Schedule) error {
	if a.scheduler != nil {
		return fmt.Errorf("illegal state: scheduler is already running")
	}
	a.model = model
	a.scheduler = &controller.Scheduler{}
	a.scheduler.SetLogger(a.Log)
	err := a.scheduler.Start(model)
	if err == nil {
		a.SendEvent("config", model)
	}
	return err
}

// Stop the scheduler module.
func (a *SchedulerApp) Stop() error {
	var err error
	if a.scheduler != nil {
		tmp := a.scheduler
		a.scheduler = nil
		err = tmp.Stop()
	}
	return err
}

// Schedule a new task or re-schedules and existing task.
func (a *SchedulerApp) Schedule(task *model.Task) (*string, error) {
	if a.scheduler != nil {
		err := a.Cancel(task.Uuid)
		if err != nil {
			a.Log.Warningf("cancel failed %s", err)
		}
		err = a.scheduler.Schedule(task)
		if err == nil {
			a.model.Tasks = append(a.model.Tasks, task)
			a.SendEvent("config", a.model)
		}
		copy := task.Uuid
		return &copy, err
	}
	return nil, fmt.Errorf("cannot schedule a task while the scheduler is stopped")
}

// Cancel an existing task.
func (a *SchedulerApp) Cancel(taskID string) error {
	if a.scheduler != nil {
		var err error
		found := -1
		for i, t := range a.model.Tasks {
			if t.Uuid == taskID {
				found = i
				break
			}
		}
		if found > -1 {
			err = a.scheduler.Cancel(taskID)
			if err == nil {
				a.model.Tasks = append(a.model.Tasks[0:found], a.model.Tasks[found+1:]...)
				a.SendEvent("config", a.model)
			}
		} else {
			err = nil
		}
		return err
	}
	return fmt.Errorf("cannot cancel a task while the scheduler is stopped")
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

	topic := fmt.Sprintf("$node/%s/app/%s/service/%s", config.Serial(), app.Info.ID, "scheduler")
	announcement := &nmodel.ServiceAnnouncement{
		Schema: "http://schema.ninjablocks.com/service/scheduler",
	}

	service, err := app.Conn.ExportService(app, topic, announcement)
	if err != nil {
		app.Log.Fatalf("failed to export scheduler service: %v", err)
	}
	_ = service

	support.WaitUntilSignal()
}
