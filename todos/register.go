package todos

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/middleware"
	"github.com/skaisanlahti/try-go-htmx/todos/handlers"
	"github.com/skaisanlahti/try-go-htmx/todos/psql"
)

func RegisterHandlers(router *http.ServeMux, database *sql.DB, todoPage *template.Template) {
	repository := psql.NewTodoRepository(database)
	getTodoPage := handlers.NewGetTodoPageHandler(repository, handlers.NewTemplateGetTodoPageView(todoPage))
	getTodoList := handlers.NewGetTodoListHandler(repository, handlers.NewGetTodoListView(todoPage))
	addTodo := handlers.NewAddTodoHandler(repository, handlers.NewTemplateAddTodoView(todoPage))
	toggleTodo := handlers.NewToggleTodoHandler(repository, handlers.NewTemplateToggleTodoView(todoPage))
	removeTodo := handlers.NewRemoveTodoHandler(repository)

	router.Handle("/todos/remove", middleware.Log(removeTodo))
	router.Handle("/todos/toggle", middleware.Log(toggleTodo))
	router.Handle("/todos/add", middleware.Log(addTodo))
	router.Handle("/todos/list", middleware.Log(getTodoList))
	router.Handle("/todos", middleware.Log(getTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
