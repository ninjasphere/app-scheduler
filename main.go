package main

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/controller"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/rpc"
	"github.com/ninjasphere/go-ninja/support"
	"time"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type SchedulerApp struct {
	support.AppSupport
	service *SchedulerService
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
	return err
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

//SchedulerService is a facade used to prevent the app event handler being replaced
//by the service event handler.
type SchedulerService struct {
	scheduler   *controller.Scheduler
	log         *logger.Logger
	conn        *ninja.Connection
	model       *model.Schedule
	service     *rpc.ExportedService
	configStore func(m *model.Schedule)
}

func (s *SchedulerService) init(moduleID string) error {
	s.scheduler = &controller.Scheduler{}
	s.scheduler.SetLogger(s.log)
	s.scheduler.SetConnection(s.conn, time.Millisecond*time.Duration(config.Int(10000, "scheduler", "timeout")))
	if err := s.scheduler.Start(s.model); err != nil {
		return err
	}
	var err error
	topic := fmt.Sprintf("$node/%s/app/%s/service/%s", config.Serial(), moduleID, "scheduler")
	announcement := &nmodel.ServiceAnnouncement{
		Schema: "http://schema.ninjablocks.com/service/scheduler",
	}
	if s.service, err = s.conn.ExportService(s, topic, announcement); err != nil {
		return err
	}
	return nil
}

// Schedule a new task or re-schedules an existing task.
func (s *SchedulerService) Schedule(task *model.Task) (*string, error) {
	if s.scheduler != nil {
		err := s.Cancel(task.ID)
		if err != nil {
			s.log.Warningf("cancel failed %s", err)
		}
		err = s.scheduler.Schedule(task)
		if err == nil {
			s.model.Tasks = append(s.model.Tasks, task)
			s.configStore(s.model)
		}
		copy := task.ID
		return &copy, err
	}
	return nil, fmt.Errorf("cannot schedule a task while the scheduler is stopped")
}

// Cancel an existing task.
func (s *SchedulerService) Cancel(taskID string) error {
	if s.scheduler != nil {
		var err error
		found := -1
		for i, t := range s.model.Tasks {
			if t.ID == taskID {
				found = i
				break
			}
		}
		if found > -1 {
			err = s.scheduler.Cancel(taskID)
			if err == nil {
				s.model.Tasks = append(s.model.Tasks[0:found], s.model.Tasks[found+1:]...)
				s.configStore(s.model)
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

	support.WaitUntilSignal()
}
