package main

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type SchedulerApp struct {
	support.AppSupport
	service   *SchedulerService
	userAgent *userAgentListener
}

// Start is called after the ExportApp call is complete.
func (a *SchedulerApp) Start(m *model.Schedule) error {
	if a.service != nil {
		return fmt.Errorf("illegal state: scheduler is already running")
	}
	a.service = &SchedulerService{
		log:   a.Log,
		conn:  a.Conn,
		model: m,
		configStore: func(m *model.Schedule) {
			a.SendEvent("config", m)
		},
	}
	err := a.service.init(a.Info.ID)
	if err != nil {
		return err
	}
	a.userAgent = &userAgentListener{}
	err = a.userAgent.init(a.service)
	if err != nil {
		return err
	}
	return nil
}

// Stop the scheduler module.
func (a *SchedulerApp) Stop() error {
	var err error
	if a.service != nil {
		tmp := a.service
		a.service = nil
		err = tmp.scheduler.Stop()
	}
	return err
}

type userAgentListener struct {
	service *SchedulerService
	client  *ninja.ServiceClient
}

func (l *userAgentListener) init(service *SchedulerService) error {
	l.service = service
	l.client = l.service.conn.GetServiceClient("$device/:deviceId/channel/user-agent")
	l.client.OnEvent("schedule-task", l.OnTaskSchedule)
	l.client.OnEvent("cancel-task", l.OnTaskCancel)
	return nil
}

func (l *userAgentListener) OnTaskSchedule(m *model.Task, keys map[string]string) bool {
	_, err := l.service.Schedule(m)
	if err != nil {
		l.service.log.Errorf("failed while scheduling task received via user-agent notification from %s: %v: %v", keys["deviceId"], m, err)
	}
	return true
}

func (l *userAgentListener) OnTaskCancel(pTaskID *string, keys map[string]string) bool {
	if pTaskID == nil {
		l.service.log.Errorf("illegal argument: pTaskID is nil")
		return true
	}
	err := l.service.Cancel(*pTaskID)
	if err != nil {
		l.service.log.Errorf("failed while canceling task received from user-agent notification from %s: %s: %v", keys["deviceId"], *pTaskID, err)
	}
	return true
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
