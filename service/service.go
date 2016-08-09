package service

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/controller"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/rpc"
	"github.com/pborman/uuid"
	"time"
)

//SchedulerService is a facade used to prevent the app event handler being replaced
//by the service event handler.
type SchedulerService struct {
	Scheduler   *controller.Scheduler
	Log         *logger.Logger
	Conn        *ninja.Connection
	Model       *model.Schedule
	Service     *rpc.ExportedService
	ConfigStore func(m *model.Schedule)
}

func (s *SchedulerService) Init(moduleID string) error {
	siteID := config.MustString("siteId")

	s.Scheduler = &controller.Scheduler{}
	s.Scheduler.SetLogger(s.Log)
	s.Scheduler.SetSiteID(siteID)
	s.Scheduler.SetConfigStore(s.ConfigStore)
	s.Scheduler.SetConnection(s.Conn, time.Millisecond*time.Duration(config.Int(10000, "scheduler", "timeout")))

	var err error
	topic := fmt.Sprintf("$site/%s/service/%s", siteID, "scheduler")
	announcement := &nmodel.ServiceAnnouncement{
		Schema: "http://schema.ninjablocks.com/service/scheduler",
	}
	if s.Service, err = s.Conn.ExportService(s, topic, announcement); err != nil {
		return err
	}
	if err := s.Scheduler.Start(s.Model); err != nil {
		return err
	}
	return nil
}

// Schedule a new task or re-schedules an existing task.
func (s *SchedulerService) Schedule(task *model.Task) (*string, error) {
	if s.Scheduler != nil {
		var err error
		if task.ID != "" {
			_, err = s.Cancel(task.ID)
			if err != nil {
				// might be ok
			}
		} else {
			task.ID = uuid.NewUUID().String()
		}
		err = s.Scheduler.Schedule(task)
		var copy string
		if err == nil {
			s.Scheduler.FlushModel()
			copy = task.ID
		}
		return &copy, err
	}
	return nil, fmt.Errorf("cannot schedule a task while the scheduler is stopped")
}

// Cancel an existing task.
func (s *SchedulerService) Cancel(taskID string) (*model.Task, error) {
	if s.Scheduler != nil {
		var err error
		var m *model.Task
		if m, err = s.Scheduler.Fetch(taskID); err != nil {
			return nil, err
		} else if err = s.Scheduler.Cancel(taskID); err == nil {
			s.Scheduler.FlushModel()
		}
		return m, err
	}
	return nil, fmt.Errorf("cannot cancel a task while the scheduler is stopped")
}

// Status returns the status of a specified task.
func (s *SchedulerService) Status(taskID string) (*string, error) {
	if s.Scheduler != nil {
		status, err := s.Scheduler.Status(taskID)
		return &status, err
	}
	return nil, fmt.Errorf("cannot get the status of a task while the scheduler is stopped")
}

// Fetch the defintion of the specified task.
func (s *SchedulerService) Fetch(taskID string) (*model.Task, error) {
	if s.Scheduler != nil {
		model, err := s.Scheduler.Fetch(taskID)
		return model, err
	}
	return nil, fmt.Errorf("cannot fetch the task while the scheduler is stopped")
}

// FetchSchedule fetches the the entire schedule.
func (s *SchedulerService) FetchSchedule() (*model.Schedule, error) {
	if s.Scheduler != nil {
		return s.Model, nil
	}
	return nil, fmt.Errorf("cannot fetch the schedule while the scheduler is stopped")
}

// SetTimeZone sets the scheduler timezone. The scheduler will be restarted.
func (s *SchedulerService) SetTimeZone(timezone string) error {
	if s.Scheduler != nil {
		if err := s.Scheduler.Stop(); err != nil {
			return err
		}
		s.Model.TimeZone = timezone
		if err := s.Scheduler.Start(s.Model); err != nil {
			return err
		}
	}
	return nil
}

// SetCoordinates set the location coordinates of the schedule. The scheduler will be restarted.
func (s *SchedulerService) SetCoordinates(coordinates *model.Location) error {
	if s.Scheduler != nil {
		var err error
		if err = s.Scheduler.Stop(); err != nil {
			return err
		}
		s.Model.Location = coordinates
		if err = s.Scheduler.Start(s.Model); err != nil {
			return err
		}
	}
	return nil
}
