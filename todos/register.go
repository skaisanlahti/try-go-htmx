package todos

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/common"
)

func middleware(handler common.RouteHandler) http.Handler {
	return common.Log(common.RouteHandler(handler))
}

func RegisterHandlers(router *http.ServeMux, database *sql.DB) {
	controller := newTodoController(newSqlTodoRepository(database), newTodoView())
	router.Handle("/todos/remove", middleware(controller.removeTodo))
	router.Handle("/todos/toggle", middleware(controller.toggleTodo))
	router.Handle("/todos/add", middleware(controller.addTodo))
	router.Handle("/todos", middleware(controller.todoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
