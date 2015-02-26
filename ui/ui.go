package ui

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/app-scheduler/service"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/suit"

	nmodel "github.com/ninjasphere/go-ninja/model"
)

var log = logger.GetLogger("ui")

type ConfigService struct {
	scheduler  *service.SchedulerService
	thingModel *ninja.ServiceClient
	roomModel  *ninja.ServiceClient
	rooms      map[string]*nmodel.Room // refreshed on each request
}

func NewConfigService(scheduler *service.SchedulerService, conn *ninja.Connection) *ConfigService {
	service := &ConfigService{
		scheduler:  scheduler,
		thingModel: conn.GetServiceClient("$home/services/ThingModel"),
		roomModel:  conn.GetServiceClient("$home/services/RoomModel"),
		rooms:      make(map[string]*nmodel.Room),
	}
	return service
}

type taskForm struct {
	ID                   string   `json:"id"`
	Description          string   `json:"description"`
	OriginalDescription  string   `json:"originalDescription"`
	GeneratedDescription string   `json:"generatedDescription"`
	TurnOn               []string `json:"turnOn"`
	TurnOff              []string `json:"turnOff"`
	Time                 string   `json:"time"`
	Duration             string   `json:"duration"`
	Repeat               string   `json:"repeat"`
}

func (f *taskForm) getDBDescription() string {
	return "@ " + f.Time
}

func (f *taskForm) getUIDescription() string {

	switch f.Time {
	case "dawn", "sunrise", "dusk", "sunset":
		if f.Repeat == "daily" {
			return "@ " + f.Time + " every day"
		} else {
			return "@ " + f.Time
		}
	default:
		if f.Repeat == "daily" {
			return "@ " + f.Time + " every day"
		} else {
			if t, err := time.Parse("15:04", f.Time); err != nil {
				return "@ " + f.Time
			} else {
				if f.Duration == "" {
					t = timeToTimestamp(t.Add(time.Minute))
				} else {
					duration, _ := parseTime(f.Duration)
					t = timeToTimestamp(t)
					t = t.Add(time.Duration(duration.Hour()) * time.Hour)
					t = t.Add(time.Duration(duration.Minute()) * time.Minute)
					t = t.Add(time.Duration(duration.Second()) * time.Second)
				}
				if t.Format("2006-01-02") == time.Now().Format("2006-01-02") {
					return "@ " + f.Time + " today"
				} else {
					return "@ " + f.Time + " tomorrow"
				}
			}
		}
	}

}

type thingModel struct {
	ID       string
	Name     string
	Location *string
	On       bool
}

// convert a time to the next local timestamp greater than the current time
func timeToTimestamp(hhmmss time.Time) time.Time {
	now := time.Now()
	parsed, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+hhmmss.Format("15:04:05"), now.Location())
	if parsed.Sub(now) < 0 {
		parsed = parsed.AddDate(0, 0, 1)
	}
	return parsed
}

// transforms a task form into a task model
func toModelTask(f *taskForm) (*model.Task, error) {

	after := &model.Event{}
	before := &model.Event{}

	if f.Duration == "" {
		before.Rule = "delay"
		before.Param = "00:01:00"
	} else {
		before.Rule = "delay"
		parsed, err := parseTime(f.Duration)
		if err != nil {
			return nil, fmt.Errorf("Duration must be specified in the form hh:mm or hh:mm:ss")
		}
		before.Param = parsed.Format("15:04:05")
	}

	switch f.Time {
	case "dawn", "dusk", "sunset", "sunrise":
		f.Repeat = "daily"

	default:
		parsed, err := parseTime(f.Time)
		if err != nil {
			return nil, fmt.Errorf("At must be entered as hh:mm.")
		}

		switch f.Repeat {
		case "once":
			after.Rule = "timestamp"
			after.Param = timeToTimestamp(parsed).Format("2006-01-02 15:04:05")
		case "daily":
			after.Rule = "time-of-day"
			after.Param = parsed.Format("15:04:05")
		default:
			return nil, fmt.Errorf("repeat is not valid")
		}
	}

	openActions := []*model.Action{}
	closeActions := []*model.Action{}

	for _, a := range f.TurnOn {
		o := &model.Action{
			ActionType: "thing-action",
			Action:     "turnOn",
			SubjectID:  a,
		}
		c := &model.Action{
			ActionType: "thing-action",
			Action:     "turnOff",
			SubjectID:  a,
		}
		openActions = append(openActions, o)
		closeActions = append(closeActions, c)
	}

	for _, a := range f.TurnOff {
		o := &model.Action{
			ActionType: "thing-action",
			Action:     "turnOff",
			SubjectID:  a,
		}
		c := &model.Action{
			ActionType: "thing-action",
			Action:     "turnOn",
			SubjectID:  a,
		}
		openActions = append(openActions, o)
		closeActions = append(closeActions, c)
	}

	if f.Duration == "" {
		closeActions = []*model.Action{}
	}

	description := ""
	generateable := (f.Description == f.OriginalDescription && f.GeneratedDescription == "true") || f.Description == ""
	if generateable {
		description = f.getDBDescription()
	} else {
		description = f.Description
	}

	return &model.Task{
		ID:          f.ID,
		Description: description,
		Tags: []string{
			"config-ui",
		},
		Open:  openActions,
		Close: closeActions,
		Window: &model.Window{
			After:  after,
			Before: before,
		},
	}, nil
}

// transforms a task model into a task form
func toTaskForm(m *model.Task) (*taskForm, error) {

	if indexOf(m.Tags, "simple-ui") == 0 {
		return nil, fmt.Errorf("missing a compatible tag")
	}

	f := &taskForm{
		ID:          m.ID,
		Description: m.Description,
		TurnOn:      make([]string, 0),
		TurnOff:     make([]string, 0),
	}

	switch m.Window.After.Rule {
	case "time-of-day":
		if parsed, err := time.Parse("15:04:05", m.Window.After.Param); err != nil {
			return nil, fmt.Errorf("invalid window.after.param: %v", err)
		} else {
			if parsed.Second() == 0 {
				f.Time = parsed.Format("15:04")
			} else {
				f.Time = parsed.Format("15:04:05")
			}
			f.Repeat = "daily"
		}
	case "timestamp":
		if parsed, err := time.Parse("2006-01-02 15:04:05", m.Window.After.Param); err != nil {
			return nil, fmt.Errorf("invalid window.after.param: %v", err)
		} else {
			if parsed.Second() == 0 {
				f.Time = parsed.Format("15:04")
			} else {
				f.Time = parsed.Format("15:04:05")
			}
			f.Repeat = "once"
		}
	default:
		return nil, fmt.Errorf("invalid after rule: %v", m.Window.After.Rule)
	}

	switch m.Window.Before.Rule {
	case "delay":
		if len(m.Close) > 0 {
			if parsed, err := time.Parse("15:04:05", m.Window.Before.Param); err != nil {
				return nil, fmt.Errorf("invalid window.before.param: %v", err)
			} else {
				f.Duration = parsed.Format("15:04:05")
			}
		} else {
			f.Duration = ""
		}
	default:
		return nil, fmt.Errorf("invalid before rule: %v", m.Window.Before.Rule)
	}

	f.GeneratedDescription = "false"
	f.Description = m.Description

	dbDesc := f.getDBDescription()
	if m.Description == dbDesc {
		f.GeneratedDescription = "true"
		uiDesc := f.getUIDescription()
		f.Description = uiDesc
	}

	f.OriginalDescription = f.Description

	return f, nil
}

func (c *ConfigService) error(message string) (*suit.ConfigurationScreen, error) {

	return &suit.ConfigurationScreen{
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.Alert{
						Title:        "Error",
						Subtitle:     message,
						DisplayClass: "danger",
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.ReplyAction{
				Label: "Cancel",
				Name:  "list",
			},
		},
	}, nil
}
func (c *ConfigService) list() (*suit.ConfigurationScreen, error) {

	var tasks []suit.ActionListOption

	schedule, err := c.scheduler.FetchSchedule()

	if err != nil {
		return c.error(fmt.Sprintf("Could not fetch schedule: %s", err))
	}

	c.refreshRooms()

	for _, t := range schedule.Tasks {
		if f, err := toTaskForm(t); err != nil {
			log.Debugf("skipped task (%s) because it cannot be edited: %v", t.ID, err)
			continue
		} else {
			tasks = append(tasks, suit.ActionListOption{
				Title: f.Description,
				Value: t.ID,
			})
		}
	}

	screen := suit.ConfigurationScreen{
		Title:       "Scheduler",
		DisplayIcon: "clock-o",
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.ActionList{
						Name:    "id",
						Options: tasks,
						PrimaryAction: &suit.ReplyAction{
							Name:        "edit",
							DisplayIcon: "pencil",
						},
						SecondaryAction: &suit.ReplyAction{
							Name:         "delete",
							Label:        "Delete",
							DisplayIcon:  "trash",
							DisplayClass: "danger",
						},
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.CloseAction{
				Label: "Close",
			},
			suit.ReplyAction{
				Label:        "New Task",
				Name:         "new",
				DisplayClass: "success",
				DisplayIcon:  "star",
			},
		},
	}

	return &screen, nil
}

func (c *ConfigService) Configure(request *nmodel.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	log.Infof("Incoming configuration request. Action:%s Data:%s", request.Action, string(request.Data))

	switch request.Action {
	case "list":
		fallthrough
	case "":
		return c.list()
	case "new":
		return c.edit(&model.Task{
			Window: &model.Window{
				After:  &model.Event{},
				Before: &model.Event{},
			},
		})
	case "edit":

		var vals map[string]string
		json.Unmarshal(request.Data, &vals)
		task, err := c.scheduler.Fetch(vals["id"])

		if err != nil {
			return c.error(fmt.Sprintf("Could not find task with id %s : %s", vals["task"], err))
		}

		return c.edit(task)
	case "delete":

		var vals map[string]string
		json.Unmarshal(request.Data, &vals)
		_, err := c.scheduler.Cancel(vals["id"])

		if err != nil {
			return nil, fmt.Errorf("Failed to delete task %s: %s", vals["id"], err)
		}

		return c.list()
	case "save":

		form := &taskForm{}
		if err := json.Unmarshal(request.Data, form); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal save task request %s: %s", request.Data, err)
		}
		task, err := toModelTask(form)
		if err != nil {
			return nil, fmt.Errorf("Failed to transform form into task: %v", err)
		}
		task = task.Migrate()
		if _, err := c.scheduler.Schedule(task); err != nil {
			return nil, fmt.Errorf("Failed to save task %s: %s", request.Data, err)
		}

		return c.list()
	default:
		return c.error(fmt.Sprintf("Unknown action: %s", request.Action))
	}
}

func (c *ConfigService) edit(task *model.Task) (*suit.ConfigurationScreen, error) {

	var form *taskForm
	var err error

	c.refreshRooms()

	if form, err = toTaskForm(task); err != nil {
		return c.error(fmt.Sprintf("Could not load form from model %s", err))
	}

	onOffThings, err := c.getOnOffThings()

	if err != nil {
		return c.error(fmt.Sprintf("Could not fetch all things: %s", err))
	}

	var turnOnOptions []suit.OptionGroupOption
	for _, s := range onOffThings {
		title := s.Name
		if s.Location != nil {
			if room, ok := c.rooms[*s.Location]; ok {
				title = fmt.Sprintf("%s in %s", s.Name, room.Name)
			}
		}
		turnOnOptions = append(turnOnOptions, suit.OptionGroupOption{
			Title:    title,
			Value:    s.ID,
			Selected: containsThingAction(task, "turnOn", s.ID),
		})
	}

	var turnOffOptions []suit.OptionGroupOption
	for _, s := range onOffThings {
		title := s.Name
		if s.Location != nil {
			if room, ok := c.rooms[*s.Location]; ok {
				title = fmt.Sprintf("%s in %s", s.Name, room.Name)
			}
		}
		turnOffOptions = append(turnOffOptions, suit.OptionGroupOption{
			Title:    title,
			Value:    s.ID,
			Selected: containsThingAction(task, "turnOff", s.ID),
		})
	}

	title := "New Scheduled Task"
	if task.ID != "" {
		title = "Edit Scheduled Task"
	}

	screen := suit.ConfigurationScreen{
		Title:       title,
		DisplayIcon: "clock-o",
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.InputHidden{
						Name:  "id",
						Value: form.ID,
					},
					suit.InputHidden{
						Name:  "originalDescription",
						Value: form.Description,
					},
					suit.InputHidden{
						Name:  "generatedDescription",
						Value: fmt.Sprintf("%v", form.GeneratedDescription),
					},
					suit.Separator{},
					suit.OptionGroup{
						Name:    "turnOn",
						Title:   "Turn on",
						Options: turnOnOptions,
					},
					suit.OptionGroup{
						Name:    "turnOff",
						Title:   "Turn off",
						Options: turnOffOptions,
					},
					suit.Separator{},
					suit.InputText{
						Name:        "time",
						Title:       "At",
						Placeholder: "hh:mm|dawn|sunrise|sunset|dusk",
						Value:       form.Time,
					},
					suit.InputText{
						Title:       "Duration",
						Name:        "duration",
						Placeholder: "hh:mm:ss",
						Value:       form.Duration,
					},
					suit.InputText{
						Title:       "Name",
						Name:        "description",
						Placeholder: "My Task",
						Value:       form.Description,
					},
					suit.RadioGroup{
						Name:  "repeat",
						Title: "Repeat",
						Value: form.Repeat,
						Options: []suit.RadioGroupOption{
							suit.RadioGroupOption{
								Title:       "Once",
								Value:       "once",
								DisplayIcon: "bolt",
							},
							suit.RadioGroupOption{
								Title:       "Daily",
								Value:       "daily",
								DisplayIcon: "repeat",
							},
						},
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.ReplyAction{
				Label: "Cancel",
				Name:  "list",
			},
			suit.ReplyAction{
				Label:        "Save",
				Name:         "save",
				DisplayClass: "success",
				DisplayIcon:  "star",
			},
		},
	}

	return &screen, nil
}

func (c *ConfigService) getOnOffThings() ([]*thingModel, error) {

	var things []*nmodel.Thing

	err := c.thingModel.Call("fetchAll", []interface{}{}, &things, time.Second*20)
	//err = client.Call("fetch", "c7ac05e0-9999-4d93-bfe3-a0b4bb5e7e78", &thing)

	if err != nil {
		return nil, fmt.Errorf("Failed to get things!: %s", err)
	}

	onOffThings := []*thingModel{}

	for _, thing := range things {
		if channels := thing.Device.GetChannelsByProtocol("on-off"); len(channels) > 0 && thing.Promoted {
			on := false
			ch := channels[0]
			if ch.LastState != nil {
				if state, ok := ch.LastState.(map[string]interface{}); ok {
					if payload, ok := state["payload"]; ok {
						on, _ = payload.(bool)
					}
				}
			}
			m := &thingModel{
				ID:       thing.ID,
				Location: thing.Location,
				Name:     thing.Name,
				On:       on,
			}
			onOffThings = append(onOffThings, m)
		}
	}

	return onOffThings, nil
}

func (c *ConfigService) refreshRooms() error {

	var rooms []*nmodel.Room

	err := c.roomModel.Call("fetchAll", []interface{}{}, &rooms, time.Second*20)
	//err = client.Call("fetch", "c7ac05e0-9999-4d93-bfe3-a0b4bb5e7e78", &thing)

	if err != nil {
		return fmt.Errorf("Failed to get rooms!: %s", err)
	}

	result := make(map[string]*nmodel.Room)
	for _, r := range rooms {
		result[r.ID] = r
	}

	c.rooms = result
	return nil
}

func i(i int) *int {
	return &i
}

func indexOf(s []string, m string) int {
	for i, e := range s {
		if e == m {
			return i
		}
	}
	return -1
}

func containsThingAction(task *model.Task, action, thingID string) bool {
	for _, a := range task.Open {
		if a.SubjectID == "thing:"+thingID && a.ActionType == "thing-action" && a.ActionType == action {
			return true
		}
	}

	return false
}

func parseTime(in string) (time.Time, error) {
	var (
		parsed time.Time
		err    error
	)
	parsed, err = time.Parse("15:04:05", in)
	if err != nil {
		parsed, err = time.Parse("15:04", in)
	}
	return parsed, err
}
