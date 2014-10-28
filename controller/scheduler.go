package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/logger"
)

var log logger.Logger

type scheduler struct {
	model *model.Schedule
	tasks map[string]*task
}

func (s *scheduler) start(m *model.Schedule) error {
	s.model = m
	s.tasks = make(map[string]*task)
	for _, t := range m.Tasks {
		task := &task{}
		err := task.init(t)
		if err == nil {
			s.tasks[t.Uuid] = task
			go task.loop()
		} else {
			log.Errorf("failed to initialize task id '%s' because %s", t.Uuid, err)
		}
	}
	return nil
}

func (s *scheduler) stop() error {
	for _, t := range s.tasks {
		t.quit <- struct{}{}
	}
	return nil
}
