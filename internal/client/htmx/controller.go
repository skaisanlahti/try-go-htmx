package htmx

import (
	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type controller struct {
	model *model
	view  *view
}

func newController(security *security.SecurityService, todo *todo.TodoService) *controller {
	return &controller{newModel(security, todo), newView()}
}
