package htmx

import (
	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type applicationFacade struct {
	*security.SecurityService
	*todo.TodoService
}

func newApplicationFacade(security *security.SecurityService, todo *todo.TodoService) *applicationFacade {
	return &applicationFacade{security, todo}
}
