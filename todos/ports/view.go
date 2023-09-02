package ports

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/models"
)

type View interface {
	RenderPage(response http.ResponseWriter, todoPage models.TodoPage) error
	RenderForm(response http.ResponseWriter, todoPage models.TodoPage) error
	RenderList(response http.ResponseWriter, todoPage models.TodoPage) error
	RenderItem(response http.ResponseWriter, todo models.Todo) error
}
