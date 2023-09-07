package todos

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/middleware"
	"github.com/skaisanlahti/try-go-htmx/todos/htmx"
	"github.com/skaisanlahti/try-go-htmx/todos/psql"
	"github.com/skaisanlahti/try-go-htmx/todos/templates"
)

func MapHtmxHandlers(router *http.ServeMux, database *sql.DB) {
	repository := psql.NewTodoRepository(database)
	todoPage := templates.ParseTemplates().TodoPage
	removeTodo := htmx.NewRemoveTodoHandler(repository)
	toggleTodo := htmx.NewToggleTodoHandler(repository, htmx.NewHtmxToggleTodoView(todoPage))
	addTodo := htmx.NewAddTodoHandler(repository, htmx.NewHtmxAddTodoView(todoPage))
	getTodoList := htmx.NewGetTodoListHandler(repository, htmx.NewHtmxGetTodoListView(todoPage))
	getTodoPage := htmx.NewGetTodoPageHandler(repository, htmx.NewHtmxGetTodoPageView(todoPage))

	router.Handle("/todos/remove", middleware.Log(removeTodo))
	router.Handle("/todos/toggle", middleware.Log(toggleTodo))
	router.Handle("/todos/add", middleware.Log(addTodo))
	router.Handle("/todos/list", middleware.Log(getTodoList))
	router.Handle("/todos", middleware.Log(getTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
