package htmx

import (
	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type model struct {
	*security.SecurityService
	*todo.TodoService
}

func newModel(security *security.SecurityService, todo *todo.TodoService) *model {
	return &model{security, todo}
}
