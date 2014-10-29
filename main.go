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

type SchedulerApp struct {
	support.AppSupport
	scheduler *controller.Scheduler
	model     *model.Schedule
}

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

func (a *SchedulerApp) Stop() error {
	var err error
	if a.scheduler != nil {
		tmp := a.scheduler
		a.scheduler = nil
		err = tmp.Stop()
	}
	return err
}

func (a *SchedulerApp) Schedule(task *model.Task) (string, error) {
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
		return task.Uuid, err
	} else {
		return "", fmt.Errorf("cannot schedule a task while the scheduler is stopped")
	}
}

func (a *SchedulerApp) Cancel(taskId string) error {
	if a.scheduler != nil {
		var err error
		found := -1
		for i, t := range a.model.Tasks {
			if t.Uuid == taskId {
				found = i
				break
			}
		}
		if found > -1 {
			err = a.scheduler.Cancel(taskId)
			if err == nil {
				a.model.Tasks = append(a.model.Tasks[0:found], a.model.Tasks[found+1:]...)
				a.SendEvent("config", a.model)
			}
		} else {
			err = nil
		}
		return err
	} else {
		return fmt.Errorf("cannot cancel a task while the scheduler is stopped")
	}
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
