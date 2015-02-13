package controller

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/config"
)

type presetsAction struct {
	baseAction
}

func (a *presetsAction) actuate(ctx *actuationContext) error {
	siteID := config.MustString("siteId")
	topic := fmt.Sprintf("$site/%s/service/%s", siteID, "presets")
	client := ctx.conn.GetServiceClient(topic)
	id := a.getModel().GetSceneID()
	if id != "" {
		params := []string{id}
		reply := &struct{}{}
		return client.Call(a.model.Action, params, reply, ctx.timeout)
	} else {
		return fmt.Errorf("The scene id for presets action was empty. The actuation did nothing.")
	}
}
