package core

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
)

func (this *Service) GetTodoPage(response http.ResponseWriter, request *http.Request) error {
	todos, err := this.Database.GetTodos()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	data := models.NewTodoPage()
	data.Todos = todos
	return this.View.RenderPage(response, data)
}
