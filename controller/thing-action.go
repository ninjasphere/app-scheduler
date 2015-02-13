package controller

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/model"
)

type thingAction struct {
	baseAction
}

func (a *thingAction) actuate(ctx *actuationContext) error {
	conn := ctx.conn
	client := ctx.thingClient
	timeout := ctx.timeout

	// acquire the model
	thing := &model.Thing{}
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
