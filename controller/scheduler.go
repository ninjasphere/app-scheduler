package controller

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"time"
)

var log = logger.GetLogger("")

type startRequest struct {
	model       *model.Task
	reply       chan error
	updateModel bool
}

type cancelRequest struct {
	id    string
	reply chan error
}

type actuationRequest struct {
	action *action
	reply  chan error
}

type statusRequest struct {
	taskID string
	status string
	err    error
	reply  chan *statusRequest
}

// Scheduler is a controller that coordinates the execution of the tasks specified by a model schedule.
type Scheduler struct {
	conn        *ninja.Connection
	thingClient *ninja.ServiceClient
	siteClient  *ninja.ServiceClient
	timeout     time.Duration
	dirty       bool
	siteID      string
	model       *model.Schedule
	started     map[string]*task
	stopped     chan struct{}
	shutdown    chan bool
	tasks       chan startRequest
	cancels     chan cancelRequest
	actuations  chan actuationRequest
	status      chan *statusRequest
	flush       chan struct{}
	configStore func(m *model.Schedule)
}

func (s *Scheduler) flushModel() {
	if s.dirty {
		s.configStore(s.model)
		s.dirty = false
	}
}

// The control loop of the scheduler. It is responsible for admitting
// new tasks, reaping completed tasks, cancelling running tasks.
func (s *Scheduler) loop() {
	var quit = false

	reap := make(chan string)

	for !quit || len(s.started) > 0 {
		select {
		case quit = <-s.shutdown:
			for taskID, t := range s.started {
				log.Debugf("signaled %s", taskID)
				t.quit <- struct{}{}
			}

		case _ = <-s.flush:
			s.flushModel()

		case taskID := <-reap:
			log.Debugf("reaped %s", taskID)
			delete(s.started, taskID)

		case startReq := <-s.tasks:
			taskID := startReq.model.ID
			runner := &task{}
			err := runner.init(startReq.model, s.actuations)
			if err == nil {
				s.started[taskID] = runner
				if startReq.updateModel {
					s.dirty = true
					s.model.Tasks = append(s.model.Tasks, startReq.model)
				}
				go func() {
					defer func() {
						log.Debugf("exiting %s", taskID)
						reap <- taskID
					}()
					permanentlyClosed := runner.loop()
					if permanentlyClosed {
						reply := make(chan error, 1)
						s.cancels <- cancelRequest{taskID, reply}
						_ = <-reply
						s.flush <- struct{}{}
					}
				}()
				log.Debugf("started %s", taskID)
			}
			startReq.reply <- err

		case cancelReq := <-s.cancels:
			var err error

			var found = -1
			for i, m := range s.model.Tasks {
				if m.ID == cancelReq.id {
					found = i
					break
				}
			}

			if found >= 0 {
				s.model.Tasks = append(s.model.Tasks[0:found], s.model.Tasks[found+1:]...)
				s.dirty = true
			}

			if runner, ok := s.started[cancelReq.id]; ok {
				if found < 0 {
					err = fmt.Errorf("found task %s in runtime but not in model", cancelReq.id)
				}
				runner.quit <- struct{}{}
			} else {
				if found >= 0 {
					err = fmt.Errorf("found task %s in model but not in runtime", cancelReq.id)
				}
			}

			cancelReq.reply <- err

		case actuationReq := <-s.actuations:
			err := actuationReq.action.actuate(s.conn, s.thingClient, s.timeout)
			actuationReq.reply <- err
		case statusReq := <-s.status:
			if t, ok := s.started[statusReq.taskID]; ok {
				statusReq.status = t.status
			} else {
				statusReq.err = fmt.Errorf("Task %s not found", statusReq.taskID)
				statusReq.status = statusReq.err.Error()
			}
			statusReq.reply <- statusReq

		}

	}

	s.stopped <- struct{}{}

}

// Start the scheduler. Iterate over the model schedule, creating and starting tasks for each Task model found.
func (s *Scheduler) Start(m *model.Schedule) error {
	s.model = m
	s.dirty = true
	s.shutdown = make(chan bool)
	s.started = make(map[string]*task)
	s.stopped = make(chan struct{})
	s.tasks = make(chan startRequest)
	s.cancels = make(chan cancelRequest)
	s.actuations = make(chan actuationRequest)
	s.status = make(chan *statusRequest)
	s.flush = make(chan struct{})

	var err error

	// update schedule model timezon and location paramters from the current site parameters provided that they
	// are not nil or empty.

	siteModel := &nmodel.Site{}
	if err = s.siteClient.Call("fetch", s.siteID, siteModel, s.timeout); err != nil {
		return fmt.Errorf("error while retrieving model site: %v", err)
	}

	if siteModel.TimeZoneID != nil && *siteModel.TimeZoneID != "" {
		s.model.TimeZone = *siteModel.TimeZoneID
	}

	if siteModel.Latitude != nil && siteModel.Longitude != nil {
		s.model.Location = &model.Location{}
		s.model.Location.Latitude = *siteModel.Latitude
		s.model.Location.Longitude = *siteModel.Longitude
		s.model.Location.Altitude = 0.0
	}

	var loc *time.Location
	// set the timezone of the clock
	if loc, err = time.LoadLocation(s.model.TimeZone); err != nil {
		log.Warningf("error while setting timezone (%s): %s", s.model.TimeZone, err)
		loc, _ = time.LoadLocation("Local")
		s.model.TimeZone = "Local"
	}
	clock.ResetCoordinates()
	clock.SetLocation(loc)

	// set the coordinates of the clock
	if s.model.Location != nil {
		clock.SetCoordinates(s.model.Location.Latitude, s.model.Location.Longitude, s.model.Location.Altitude)
	}

	go s.loop()

	var errors []error

	for _, t := range m.Tasks {
		reply := make(chan error)
		s.tasks <- startRequest{t, reply, false}
		err := <-reply
		if err != nil {
			errors = append(errors, err)
		}
	}

	s.flush <- struct{}{}

	if len(errors) > 1 {
		return fmt.Errorf("errors %v", errors)
	} else if len(errors) == 1 {
		return errors[0]
	} else {
		return nil
	}
}

// Stop the scheduler.
func (s *Scheduler) Stop() error {
	s.shutdown <- true
	<-s.stopped
	return nil
}

// Schedule the specified task. Starts a task controller for the specified task model.
func (s *Scheduler) Schedule(m *model.Task) error {
	reply := make(chan error)
	s.tasks <- startRequest{m, reply, true}
	err := <-reply
	return err
}

// FlushModel ensures that any pending updates to the model have been flushed back to the
// application configuration.
func (s *Scheduler) FlushModel() {
	s.flush <- struct{}{}
}

// Cancel the specified task. Stops the task controller for the specified task.
func (s *Scheduler) Cancel(taskID string) error {
	reply := make(chan error)
	s.cancels <- cancelRequest{taskID, reply}
	err := <-reply
	return err
}

// Status answers the staus of the specified task.
func (s *Scheduler) Status(taskID string) (string, error) {
	request := &statusRequest{
		taskID: taskID,
		err:    nil,
		status: "",
		reply:  make(chan *statusRequest),
	}
	s.status <- request
	_ = <-request.reply
	return request.status, request.err
}

// SetLogger sets the logger to be used by the scheduler component.
func (s *Scheduler) SetLogger(logger *logger.Logger) {
	if logger != nil {
		log = logger
	}
}

// SetConnection configure's the scheduler's connection and the default timeout
// for requests sent on the connection.
func (s *Scheduler) SetConnection(conn *ninja.Connection, timeout time.Duration) {
	s.conn = conn
	s.timeout = timeout
	s.thingClient = s.conn.GetServiceClient("$home/services/ThingModel")
	s.siteClient = s.conn.GetServiceClient("$home/services/SiteModel")
}

// SetConfigStore sets the function used to write updates to the schedule
func (s *Scheduler) SetConfigStore(store func(m *model.Schedule)) {
	s.configStore = store
}

// SetSiteID sets the site identifier of the scheduler
func (s *Scheduler) SetSiteID(id string) {
	s.siteID = id
}
