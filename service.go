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
	"time"
)

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
	siteID := config.MustString("siteId")

	s.scheduler = &controller.Scheduler{}
	s.scheduler.SetLogger(s.log)
	s.scheduler.SetSiteID(siteID)
	s.scheduler.SetConfigStore(s.configStore)
	s.scheduler.SetConnection(s.conn, time.Millisecond*time.Duration(config.Int(10000, "scheduler", "timeout")))

	var err error
	topic := fmt.Sprintf("$site/%s/service/%s", siteID, "scheduler")
	announcement := &nmodel.ServiceAnnouncement{
		Schema: "http://schema.ninjablocks.com/service/scheduler",
	}
	if s.service, err = s.conn.ExportService(s, topic, announcement); err != nil {
		return err
	}
	if err := s.scheduler.Start(s.model); err != nil {
		return err
	}
	return nil
}

// Schedule a new task or re-schedules an existing task.
func (s *SchedulerService) Schedule(task *model.Task) (*string, error) {
	if s.scheduler != nil {
		err := s.Cancel(task.ID)
		if err != nil {
			// might be ok
		}
		err = s.scheduler.Schedule(task)
		var copy string
		if err == nil {
			s.scheduler.FlushModel()
			copy = task.ID
		}
		return &copy, err
	}
	return nil, fmt.Errorf("cannot schedule a task while the scheduler is stopped")
}

// Cancel an existing task.
func (s *SchedulerService) Cancel(taskID string) error {
	if s.scheduler != nil {
		if err := s.scheduler.Cancel(taskID); err == nil {
			s.scheduler.FlushModel()
		} else {
			return err
		}
	}
	return fmt.Errorf("cannot cancel a task while the scheduler is stopped")
}

// Status returns the status of a specified task.
func (s *SchedulerService) Status(taskID string) (*string, error) {
	if s.scheduler != nil {
		status, err := s.scheduler.Status(taskID)
		return &status, err
	}
	return nil, fmt.Errorf("cannot get the status of a task while the scheduler is stopped")
}

// Fetch the defintion of the specified task.
func (s *SchedulerService) Fetch(taskID string) (*model.Task, error) {
	if s.scheduler != nil {
		model, err := s.scheduler.Fetch(taskID)
		return model, err
	}
	return nil, fmt.Errorf("cannot fetch the task while the scheduler is stopped")
}

// Fetch the entire schedule.
func (s *SchedulerService) FetchSchedule() (*model.Schedule, error) {
	if s.scheduler != nil {
		return s.model, nil
	}
	return nil, fmt.Errorf("cannot fetch the schedule while the scheduler is stopped")
}

// SetTimeZone sets the scheduler timezone. The scheduler will be restarted.
func (s *SchedulerService) SetTimeZone(timezone string) error {
	if s.scheduler != nil {
		if err := s.scheduler.Stop(); err != nil {
			return err
		}
		s.model.TimeZone = timezone
		if err := s.scheduler.Start(s.model); err != nil {
			return err
		}
	}
	return nil
}

// SetCoordinates set the location coordinates of the schedule. The scheduler will be restarted.
func (s *SchedulerService) SetCoordinates(coordinates *model.Location) error {
	if s.scheduler != nil {
		var err error
		if err = s.scheduler.Stop(); err != nil {
			return err
		}
		s.model.Location = coordinates
		if err = s.scheduler.Start(s.model); err != nil {
			return err
		}
	}
	return nil
}
