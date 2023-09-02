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
	router.Handle("/todos/remove", middleware(service.RemoveTask))
	router.Handle("/todos/toggle", middleware(service.ToggleTask))
	router.Handle("/todos/add", middleware(service.AddTask))
	router.Handle("/todos", middleware(service.DisplayPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
