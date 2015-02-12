package controller

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/go-ninja/api"
	nmodel "github.com/ninjasphere/go-ninja/model"
	"time"
)

type action interface {
	actuate(conn *ninja.Connection, client *ninja.ServiceClient, timeout time.Duration) error
	getModel() *model.Action
}

type baseAction struct {
	model *model.Action
}

type thingAction struct {
	baseAction
}

func (a *baseAction) getModel() *model.Action {
	return a.model
}

func (a *thingAction) actuate(conn *ninja.Connection, client *ninja.ServiceClient, timeout time.Duration) error {
	// acquire the model
	thing := &nmodel.Thing{}
	if err := client.Call("fetch", a.model.GetThingID(), thing, timeout); err != nil {
		return err
	}

	// iterate across all channels
	if thing.Device == nil || thing.Device.Channels == nil {
		return fmt.Errorf("'%s' does not have any channels", a.model.GetThingID())
	}

	// acquire matching topics
	topics := make([]string, 0, 0)
	for _, ch := range *thing.Device.Channels {
		if ch.ServiceAnnouncement.SupportedMethods == nil {
			continue
		}
		for _, m := range *ch.ServiceAnnouncement.SupportedMethods {
			if m == a.model.Action {
				topics = append(topics, ch.ServiceAnnouncement.Topic)
				break
			}
		}
	}

	if len(topics) == 0 {
		return fmt.Errorf("no topics supporting the '%s' method were found on '%s'", a.model.Action, a.model.GetThingID())
	}

	errors := make([]error, 0, 0)

	// call matching topics
	for _, topic := range topics {
		client := conn.GetServiceClient(topic)
		params := &struct{}{}
		reply := &struct{}{}
		err := client.Call(a.model.Action, params, reply, timeout)
		if err != nil {
			errors = append(errors, err)
		}
	}

	// collate errors
	if len(errors) > 1 {
		return fmt.Errorf("%v", errors)
	} else if len(errors) == 1 {
		return errors[0]
	}

	return nil
}

func newAction(m *model.Action) (action, error) {
	var a action

	switch m.ActionType {
	case "thing-action":
	default:
		return nil, fmt.Errorf("'%s' is not a supported action type", m.ActionType)
	}

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

	return a, nil
}
