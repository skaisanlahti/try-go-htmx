package core

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
)

func (this *Controller) AddTodo(response http.ResponseWriter, request *http.Request) error {
	task := request.FormValue("task")
	data := models.NewTodoPage()
	if task == "" {
		data.Error = "Task can't be empty."
		this.View.RenderForm(response, data)
		return nil
	}

	newTodo := models.NewTask(task)
	if err := this.Database.AddTodo(newTodo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	response.Header().Add("HX-Trigger", "GetTodoList")
	return this.View.RenderForm(response, data)
}
