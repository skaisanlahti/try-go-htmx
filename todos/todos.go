package todos

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/logging"
	"github.com/skaisanlahti/try-go-htmx/todos/handlers"
	"github.com/skaisanlahti/try-go-htmx/todos/repositories"
	"github.com/skaisanlahti/try-go-htmx/todos/templates"
	"github.com/skaisanlahti/try-go-htmx/users/sessions"
)

func UseTodoRoutes(router *http.ServeMux, database *sql.DB, store *sessions.Store) {
	repository := repositories.NewPsqlTodoRepository(database)
	todoPage := templates.ParseTemplates().TodoPage

	getTodoPage := handlers.NewGetTodoPageHandler(repository, handlers.NewHtmxGetTodoPageView(todoPage))
	getTodoList := handlers.NewGetTodoListHandler(repository, handlers.NewHtmxGetTodoListView(todoPage))
	addTodo := handlers.NewAddTodoHandler(repository, handlers.NewHtmxAddTodoView(todoPage))
	toggleTodo := handlers.NewToggleTodoHandler(repository, handlers.NewHtmxToggleTodoView(todoPage))
	removeTodo := handlers.NewRemoveTodoHandler(repository)

	router.Handle("/todos/remove", logging.LogRequest(sessions.RequireSession(removeTodo, store)))
	router.Handle("/todos/toggle", logging.LogRequest(sessions.RequireSession(toggleTodo, store)))
	router.Handle("/todos/add", logging.LogRequest(sessions.RequireSession(addTodo, store)))
	router.Handle("/todos/list", logging.LogRequest(sessions.RequireSession(getTodoList, store)))
	router.Handle("/todos", logging.LogRequest(sessions.RequireSession(getTodoPage, store)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
