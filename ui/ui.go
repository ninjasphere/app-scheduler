package ui

import (
	"encoding/json"
	"fmt"

	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/app-scheduler/service"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/suit"

	nmodel "github.com/ninjasphere/go-ninja/model"
)

var log = logger.GetLogger("ui")

type ConfigService struct {
	scheduler *service.SchedulerService
}

func NewConfigService(scheduler *service.SchedulerService) *ConfigService {
	return &ConfigService{scheduler}
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

	for _, t := range schedule.Tasks {
		tasks = append(tasks, suit.ActionListOption{
			Title: t.Description,
			Value: t.ID,
		})
	}

	screen := suit.ConfigurationScreen{
		Title: "Scheduler",
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
				Label:        "New Scheduled Task",
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
		return c.edit(&model.Task{})
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

		task := &model.Task{}
		err := json.Unmarshal(request.Data, task)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal save task request %s: %s", request.Data, err)
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

	/*var sensorOptions []suit.OptionGroupOption
	sensors, err := getSensors()
	if err != nil {
		return c.error(fmt.Sprintf("Could not find sensors: %s", err))
	}

	for _, s := range sensors {
		sensorOptions = append(sensorOptions, suit.OptionGroupOption{
			Title:    s.Name,
			Value:    s.ID,
			Selected: contains(config.Sensors, s.ID),
		})
	}

	var lightOptions []suit.OptionGroupOption
	lights, err := getLights()
	if err != nil {
		return c.error(fmt.Sprintf("Could not find lights: %s", err))
	}

	for _, s := range lights {
		lightOptions = append(lightOptions, suit.OptionGroupOption{
			Title:    s.Name,
			Value:    s.ID,
			Selected: contains(config.Lights, s.ID),
		})
	}*/

	title := "New Security Light"
	if task.ID != "" {
		title = "Edit Security Light"
	}

	screen := suit.ConfigurationScreen{
		Title: title,
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.InputHidden{
						Name:  "id",
						Value: task.ID,
					},
					suit.InputText{
						Name:        "description",
						Before:      "Name",
						Placeholder: "My Task",
						Value:       task.Description,
					},
					/*suit.OptionGroup{
						Name:           "sensors",
						Title:          "When these devices detect motion",
						MinimumChoices: 1,
						Options:        sensorOptions,
					},
					suit.OptionGroup{
						Name:           "lights",
						Title:          "Turn on these lights",
						MinimumChoices: 1,
						Options:        lightOptions,
					},
					suit.InputTimeRange{
						Name:  "time",
						Title: "When",
						Value: suit.TimeRange{
							From: config.Time.From,
							To:   config.Time.To,
						},
					},
					suit.InputText{
						Title:     "Turn off again after",
						After:     "minutes",
						Name:      "timeout",
						InputType: "number",
						Minimum:   i(0),
						Value:     config.Timeout,
					},*/
				},
			},
		},
		Actions: []suit.Typed{
			suit.CloseAction{
				Label: "Cancel",
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

func i(i int) *int {
	return &i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
