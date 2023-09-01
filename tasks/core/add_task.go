package core

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
)

func (this *Service) AddTask(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	if task == "" {
		tasks, err := this.Database.GetTasks()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return err
		}

		data := models.NewTaskPage(tasks)
		data.Error = "Task can't be empty."
		response.WriteHeader(http.StatusOK)
		this.View.RenderMain(response, data)
		return nil
	}

	newTask := models.NewTask(task)
	if err := this.Database.AddTask(newTask); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	tasks, err := this.Database.GetTasks()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data := models.NewTaskPage(tasks)
	response.WriteHeader(http.StatusCreated)
	return this.View.RenderMain(response, data)
}
