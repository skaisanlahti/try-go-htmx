package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/domain"
)

type RemoveTodoRepository interface {
	GetTodoById(id int) (domain.Todo, error)
	RemoveTodo(id int) error
}

type RemoveTodoHandler struct {
	repository RemoveTodoRepository
}

func NewRemoveTodoHandler(
	repository RemoveTodoRepository,
) *RemoveTodoHandler {
	return &RemoveTodoHandler{repository}
}

func (handler *RemoveTodoHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	id, err := extractTodoId(request.URL)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = handler.repository.GetTodoById(id)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	if err = handler.repository.RemoveTodo(id); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}
