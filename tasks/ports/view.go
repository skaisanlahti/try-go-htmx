package ports

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/tasks/models"
)

type View interface {
	RenderPage(response http.ResponseWriter, taskPage models.TaskPage) error
	RenderMain(response http.ResponseWriter, taskPage models.TaskPage) error
	RenderList(response http.ResponseWriter, taskPage models.TaskPage) error
}
