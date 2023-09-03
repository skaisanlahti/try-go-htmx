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
	controller := core.NewController(database)
	router.Handle("/todos/remove", middleware(controller.RemoveTodo))
	router.Handle("/todos/toggle", middleware(controller.ToggleTodo))
	router.Handle("/todos/add", middleware(controller.AddTodo))
	router.Handle("/todos/list", middleware(controller.GetTodoList))
	router.Handle("/todos", middleware(controller.GetTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
