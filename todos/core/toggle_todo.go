package core

import (
	"net/http"
)

func (this *Controller) ToggleTodo(response http.ResponseWriter, request *http.Request) error {
	id, err := extractTodoID(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return err
	}

	todo, err := this.Database.GetTodoByID(id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	todo.Done = !todo.Done
	if err = this.Database.UpdateTodo(todo); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return this.View.RenderItem(response, todo)
}
