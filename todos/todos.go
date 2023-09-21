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
	todoPage := templates.ParseTemplates().TodoPage
	repository := repositories.NewPsqlTodoRepository(database)

	removeTodo := handlers.NewRemoveTodoHandler(repository)
	router.Handle("/todos/remove", logging.LogRequest(sessions.RequireSession(removeTodo, store)))

	toggleTodoView := handlers.NewHtmxToggleTodoView(todoPage)
	toggleTodo := handlers.NewToggleTodoHandler(repository, toggleTodoView)
	router.Handle("/todos/toggle", logging.LogRequest(sessions.RequireSession(toggleTodo, store)))

	addTodoView := handlers.NewHtmxAddTodoView(todoPage)
	addTodo := handlers.NewAddTodoHandler(repository, addTodoView)
	router.Handle("/todos/add", logging.LogRequest(sessions.RequireSession(addTodo, store)))

	getTodoListView := handlers.NewHtmxGetTodoListView(todoPage)
	getTodoList := handlers.NewGetTodoListHandler(repository, getTodoListView)
	router.Handle("/todos/list", logging.LogRequest(sessions.RequireSession(getTodoList, store)))

	getTodoPageView := handlers.NewHtmxGetTodoPageView(todoPage)
	getTodoPage := handlers.NewGetTodoPageHandler(repository, getTodoPageView)
	router.Handle("/todos", logging.LogRequest(sessions.RequireSession(getTodoPage, store)))

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
