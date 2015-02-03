package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/app-scheduler/service"
)

type TaskRouter struct {
	scheduler *service.SchedulerService
}

func NewTaskRouter() *TaskRouter {
	return &TaskRouter{}
}

func (tr *TaskRouter) Register(r martini.Router) {
	r.Get("", tr.GetAllTasks)
	r.Post("", tr.CreateTask)
	r.Get("/:id", tr.GetTask)
	r.Put("/:id", tr.CreateTask)
	r.Delete("/:id", tr.CancelTask)
}

func writeResponse(code int, w http.ResponseWriter, response interface{}, err error) {
	if err == nil {
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf("error: %v\n", err)))
	}
}

func (tr *TaskRouter) GetAllTasks(r *http.Request, w http.ResponseWriter) {
	schedule, err := tr.scheduler.FetchSchedule()
	writeResponse(500, w, schedule, err)
}

func (tr *TaskRouter) CreateTask(r *http.Request, w http.ResponseWriter) {
	task := model.Task{}
	json.NewDecoder(r.Body).Decode(&task)
	if id, err := tr.scheduler.Schedule(&task); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("error: %v\n", err)))
	} else {
		h := w.Header()
		h.Add("Location", fmt.Sprintf("%s://%s%s/%s", "http", r.Host, r.URL.RequestURI(), *id))
		w.WriteHeader(303)
	}
}

func (tr *TaskRouter) GetTask(params martini.Params, r *http.Request, w http.ResponseWriter) {
	task, err := tr.scheduler.Fetch(params["id"])
	writeResponse(404, w, task, err)
}

func (tr *TaskRouter) CancelTask(params martini.Params, r *http.Request, w http.ResponseWriter) {
	m, err := tr.scheduler.Cancel(params["id"])
	writeResponse(404, w, m, err)
}
