package todos

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/infrastructure"
	"github.com/skaisanlahti/try-go-htmx/todos/services"
)

func middleware(handler infrastructure.RouteHandlerFunc) http.Handler {
	return infrastructure.NewLogger(infrastructure.NewErrorHandlerFunc(handler))
}

func RegisterHandlers(router *http.ServeMux, database *sql.DB, todoPageTemplate *template.Template) {
	preparer := services.NewQueryPreparer(database)
	storage := services.NewPostgreSqlStorage(preparer)
	renderer := services.NewHtmlRenderer(todoPageTemplate)
	handler := services.NewHttpHandler(storage, renderer)

	router.Handle("/todos/remove", middleware(handler.RemoveTodo))
	router.Handle("/todos/toggle", middleware(handler.ToggleTodo))
	router.Handle("/todos/add", middleware(handler.AddTodo))
	router.Handle("/todos/list", middleware(handler.GetTodoList))
	router.Handle("/todos", middleware(handler.GetTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
