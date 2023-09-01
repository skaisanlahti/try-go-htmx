package tasks

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/app"
	"github.com/skaisanlahti/try-go-htmx/tasks/core"
)

func middleware(handler app.RouteHandlerFunc) http.Handler {
	return app.NewLogger(app.NewErrorHandlerFunc(handler))
}

func RegisterHandlers(router *http.ServeMux, database *sql.DB) {
	service := core.NewService(database)
	router.Handle("/tasks/remove", middleware(service.RemoveTask))
	router.Handle("/tasks/toggle", middleware(service.ToggleTask))
	router.Handle("/tasks/add", middleware(service.AddTask))
	router.Handle("/tasks", middleware(service.DisplayPage))
}
