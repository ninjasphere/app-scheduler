package controller

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/logger"
)

var log *logger.Logger = logger.GetLogger("")

type startRequest struct {
	model *model.Task
	reply chan error
}

type cancelRequest struct {
	id    string
	reply chan error
}

type Scheduler struct {
	model    *model.Schedule
	started  map[string]*task
	stopped  chan struct{}
	shutdown chan bool
	tasks    chan startRequest
	cancels  chan cancelRequest
}

// The control loop of the scheduler. It is responsible for admitting
// new tasks, reaping completed tasks, cancelling running tasks.
func (s *Scheduler) loop() {
	var quit = false

	reap := make(chan string)

	for !quit || len(s.started) > 0 {
		select {
		case quit = <-s.shutdown:
			for taskId, t := range s.started {
				log.Debugf("signaled %s", taskId)
				t.quit <- struct{}{}
			}

		case taskId := <-reap:
			log.Debugf("reaped %s", taskId)
			delete(s.started, taskId)

		case startReq := <-s.tasks:
			taskId := startReq.model.Uuid
			runner := &task{}
			err := runner.init(startReq.model)
			if err == nil {
				s.started[taskId] = runner
				go func() {
					defer func() {
						log.Debugf("exiting %s", taskId)
						reap <- taskId
					}()
					runner.loop()
				}()
				log.Debugf("started %s", taskId)
			}
			startReq.reply <- err

		case cancelReq := <-s.cancels:
			var err error
			if runner, ok := s.started[cancelReq.id]; ok {
				runner.quit <- struct{}{}
			} else {
				err = fmt.Errorf("task id (%s) does not exist", cancelReq.id)
			}
			cancelReq.reply <- err
		}

	}

	s.stopped <- struct{}{}

}

func (s *Scheduler) Start(m *model.Schedule) error {
	s.model = m
	s.shutdown = make(chan bool)
	s.started = make(map[string]*task)
	s.stopped = make(chan struct{})
	s.tasks = make(chan startRequest)
	s.cancels = make(chan cancelRequest)

	go s.loop()

	errors := make([]error, 0)

	for _, t := range m.Tasks {
		reply := make(chan error)
		s.tasks <- startRequest{t, reply}
		err := <-reply
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 1 {
		return fmt.Errorf("errors %v", errors)
	} else if len(errors) == 1 {
		return errors[0]
	} else {
		return nil
	}
}

func (s *Scheduler) Stop() error {
	s.shutdown <- true
	close(s.shutdown)
	<-s.stopped
	return nil
}

func (s *Scheduler) Schedule(m *model.Task) error {
	reply := make(chan error)
	s.tasks <- startRequest{m, reply}
	err := <-reply
	return err
}

func (s *Scheduler) Cancel(taskId string) error {
	reply := make(chan error)
	s.cancels <- cancelRequest{taskId, reply}
	err := <-reply
	return err
}

func (s *Scheduler) SetLogger(logger *logger.Logger) {
	if logger != nil {
		log = logger
	}
}
