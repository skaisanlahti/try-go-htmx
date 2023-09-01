package core

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
)

func (this *Service) ToggleTask(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTaskID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	task, err := this.Database.GetTaskByID(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	task.Done = !task.Done
	if err = this.Database.UpdateTask(task); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	tasks, err := this.Database.GetTasks()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data := models.NewTaskPage(tasks)
	response.WriteHeader(http.StatusOK)
	return this.View.RenderList(response, data)

}
