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
	s.scheduler = &controller.Scheduler{}
	s.scheduler.SetLogger(s.log)
	s.scheduler.SetConfigStore(s.configStore)
	s.scheduler.SetConnection(s.conn, time.Millisecond*time.Duration(config.Int(10000, "scheduler", "timeout")))

	var err error
	topic := fmt.Sprintf("$node/%s/app/%s/service/%s", config.Serial(), moduleID, "scheduler")
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
			s.log.Warningf("cancel failed %s", err)
		}
		err = s.scheduler.Schedule(task)
		var copy string
		if err == nil {
			copy = task.ID
		}
		s.scheduler.FlushModel()
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
			s.scheduler.FlushModel()
		} else {
			err = nil
		}
		return err
	}
	return fmt.Errorf("cannot cancel a task while the scheduler is stopped")
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
