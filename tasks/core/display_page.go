package core

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
)

func (this *Service) DisplayPage(response http.ResponseWriter, request *http.Request) error {
	tasks, err := this.Database.GetTasks()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data := models.NewTaskPage(tasks)
	response.WriteHeader(http.StatusOK)
	return this.View.RenderPage(response, data)
}
