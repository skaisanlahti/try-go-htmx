package interfaces

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/todos/data"
)

type TemplateService interface {
	RenderPage(response http.ResponseWriter, todos []data.Todo) error
	RenderList(response http.ResponseWriter, todos []data.Todo) error
	RenderEmptyForm(response http.ResponseWriter) error
	RenderErrorForm(response http.ResponseWriter, errorMessage string) error
	RenderItem(response http.ResponseWriter, todo data.Todo) error
}
