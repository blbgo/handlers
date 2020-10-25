package taskshandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blbgo/general"
	"github.com/blbgo/httpserver"
)

type task struct {
	httpserver.RouteParams
	taskList []general.Task
	httpserver.Renderer
}

type taskView struct {
	general.Task
	ID int
}

// Setup sets up all the routes for task starting
func Setup(
	router httpserver.Router,
	rp httpserver.RouteParams,
	taskList []general.Task,
	rend httpserver.Renderer,
) {
	r := &task{RouteParams: rp, taskList: taskList, Renderer: rend}
	router.Handler("GET", "/task", http.HandlerFunc(r.task))
	router.Handler("GET", "/task/input/:id", http.HandlerFunc(r.input))
	router.Handler("POST", "/task/run", http.HandlerFunc(r.run))
}

func (r *task) task(rw http.ResponseWriter, req *http.Request) {
	r.OK(rw, "tasks", r.taskList)
}

func (r *task) input(rw http.ResponseWriter, req *http.Request) {
	params := r.Get(req)
	if len(params) != 1 {
		r.Error(rw, "error", errors.New("invalid request"))
		return
	}
	id, err := validateIDParam(params[0], r.taskList)
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	r.OK(rw, "taskInputs", &taskView{Task: r.taskList[id], ID: id})
}

func (r *task) run(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	id, err := validateIDParam(req.PostForm.Get("id"), r.taskList)
	if err != nil {
		r.Error(rw, "error", err)
		return
	}
	task := r.taskList[id]
	inputs := task.Inputs()
	var inputValues []string
	for i := range inputs {
		formValue := req.PostForm.Get(fmt.Sprint("I", i))
		inputValues = append(inputValues, formValue)
	}
	go task.Run(inputValues...)
	r.OK(rw, "task", task)
}

func validateIDParam(sid string, taskList []general.Task) (int, error) {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return 0, err
	}
	if id < 0 || id >= len(taskList) {
		return 0, errors.New("task Id out of range")
	}
	return id, nil
}
