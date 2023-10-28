package htmx

import (
	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type application struct {
	*security.SecurityService
	*todo.TodoService
}

func newApplication(security *security.SecurityService, todo *todo.TodoService) *application {
	return &application{security, todo}
}
