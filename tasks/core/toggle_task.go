package core

import (
	"net/http"
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

	response.WriteHeader(http.StatusOK)
	return this.View.RenderItem(response, task)
}
