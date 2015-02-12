package controller

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/api"
	"time"
)

type presetsAction struct {
	baseAction
}

func (a *presetsAction) actuate(conn *ninja.Connection, client *ninja.ServiceClient, timeout time.Duration) error {
	return fmt.Errorf("unimplemented function: presetsAction.actuate")
}
