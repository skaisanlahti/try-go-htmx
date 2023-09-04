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
	queryService := services.NewSqlQueryService(database)
	dataService := services.NewPostgreSqlDataService(queryService)
	htmlService := services.NewHtmlTemplateService(todoPageTemplate)
	todoHttpService := services.NewTodoHttpService(dataService, htmlService)

	router.Handle("/todos/remove", middleware(todoHttpService.RemoveTodo))
	router.Handle("/todos/toggle", middleware(todoHttpService.ToggleTodo))
	router.Handle("/todos/add", middleware(todoHttpService.AddTodo))
	router.Handle("/todos/list", middleware(todoHttpService.GetTodoList))
	router.Handle("/todos", middleware(todoHttpService.GetTodoPage))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/todos", http.StatusMovedPermanently)
	})
}
