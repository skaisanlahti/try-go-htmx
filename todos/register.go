package todos

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/app"
	"github.com/skaisanlahti/try-go-htmx/todos/core"
)

func middleware(handler app.RouteHandlerFunc) http.Handler {
	return app.NewLogger(app.NewErrorHandlerFunc(handler))
}

func RegisterHandlers(router *http.ServeMux, database *sql.DB) {
	service := core.NewService(database)
	router.Handle("/todos/remove", middleware(service.RemoveTodo))
	router.Handle("/todos/toggle", middleware(service.ToggleTodo))
	router.Handle("/todos/add", middleware(service.AddTodo))
	router.Handle("/todos/list", middleware(service.GetTodoList))
	router.Handle("/todos", middleware(service.GetTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
