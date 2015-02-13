package controller

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	"time"
)

type actuationContext struct {
	conn        *ninja.Connection
	thingClient *ninja.ServiceClient
	timeout     time.Duration
}

type action interface {
	actuate(ctx *actuationContext) error
	getModel() *model.Action
}

type baseAction struct {
	model *model.Action
}

func (a *baseAction) getModel() *model.Action {
	return a.model
}

func newAction(m *model.Action) (action, error) {
	var a action

	switch m.ActionType {
	case "thing-action":
		switch m.Action {
		case "turnOn", "turnOff", "toggle":
			a = &thingAction{
				baseAction: baseAction{
					model: m,
				},
			}
		default:
			return nil, fmt.Errorf("'%s' is an action which is not supported by the scheduler", m.Action)
		}
	case "presets-action":
		switch m.Action {
		case "apply":
			a = &presetsAction{
				baseAction: baseAction{
					model: m,
				},
			}
		default:
			return nil, fmt.Errorf("'%s' is an action which is not supported by the scheduler", m.Action)
		}
	default:
		return nil, fmt.Errorf("'%s' is not a supported action type", m.ActionType)
	}

	return a, nil
}
