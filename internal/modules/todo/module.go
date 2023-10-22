package todo

import (
	"database/sql"

	"github.com/skaisanlahti/try-go-htmx/internal/platform"
)

type todoModule struct {
	*htmxService
}

func NewModule(database *sql.DB) *todoModule {
	todoService := newTodoService(newTodoStorage(database))
	todoHtmxRenderer := newHtmxRenderer(platform.TemplateFiles)
	todoHtmxService := newHtmxService(todoService, todoHtmxRenderer)
	return &todoModule{todoHtmxService}
}
