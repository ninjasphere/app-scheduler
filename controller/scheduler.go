package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/logger"
)

var log *logger.Logger = logger.GetLogger("")

type Scheduler struct {
	model      *model.Schedule
	tasks      map[string]*task
	stopped    chan struct{}
	reaperQuit chan bool
}

func (s *Scheduler) Start(m *model.Schedule) error {
	s.model = m
	s.reaperQuit = make(chan bool)
	s.tasks = make(map[string]*task)
	s.stopped = make(chan struct{})
	reaper := make(chan string)
	for _, t := range m.Tasks {
		task := &task{}
		err := task.init(t)
		if err == nil {
			s.tasks[t.Uuid] = task
			go func() {
				defer func() {
					log.Debugf("exiting %s", t.Uuid)
					reaper <- t.Uuid
				}()
				task.loop()
			}()
		} else {
			log.Errorf("failed to initialize task id '%s' because %s", t.Uuid, err)
		}
	}

	go func() {

		var quit = false

		for !quit {
			for len(s.tasks) > 0 {
				select {
				case nonce := <-s.reaperQuit:
					_ = nonce
					quit = true
					for taskId, t := range s.tasks {
						log.Debugf("signaled %s", taskId)
						t.quit <- struct{}{}
					}
				case taskId := <-reaper:
					log.Debugf("reaped %s", taskId)
					delete(s.tasks, taskId)
				}
			}
		}

		s.stopped <- struct{}{}
	}()

	return nil
}

func (s *Scheduler) Stop() error {
	s.reaperQuit <- true
	close(s.reaperQuit)
	<-s.stopped
	return nil
}
