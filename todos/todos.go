package todos

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/middleware"
	"github.com/skaisanlahti/try-go-htmx/todos/htmx"
	"github.com/skaisanlahti/try-go-htmx/todos/psql"
	"github.com/skaisanlahti/try-go-htmx/todos/templates"
	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type SessionStore interface {
	Remove(sessionId string)
	Validate(sessionId string) (*domain.Session, bool)
	Extend(*domain.Session) *domain.Session
}

func MapHtmxHandlers(router *http.ServeMux, database *sql.DB, sessions SessionStore, mode string) {
	repository := psql.NewTodoRepository(database)
	todoPage := templates.ParseTemplates().TodoPage
	removeTodo := htmx.NewRemoveTodoHandler(repository)
	toggleTodo := htmx.NewToggleTodoHandler(repository, htmx.NewHtmxToggleTodoView(todoPage))
	addTodo := htmx.NewAddTodoHandler(repository, htmx.NewHtmxAddTodoView(todoPage))
	getTodoList := htmx.NewGetTodoListHandler(repository, htmx.NewHtmxGetTodoListView(todoPage))
	getTodoPage := htmx.NewGetTodoPageHandler(repository, htmx.NewHtmxGetTodoPageView(todoPage))

	router.Handle("/todos/remove", middleware.LogRequest(middleware.RequireSession(removeTodo, sessions, mode)))
	router.Handle("/todos/toggle", middleware.LogRequest(middleware.RequireSession(toggleTodo, sessions, mode)))
	router.Handle("/todos/add", middleware.LogRequest(middleware.RequireSession(addTodo, sessions, mode)))
	router.Handle("/todos/list", middleware.LogRequest(middleware.RequireSession(getTodoList, sessions, mode)))
	router.Handle("/todos", middleware.LogRequest(middleware.RequireSession(getTodoPage, sessions, mode)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
