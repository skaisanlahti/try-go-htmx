package todo

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/platform"
)

func MapRoutes(router *http.ServeMux, module *todoModule, factory platform.MiddlewareFactory) {
	log := factory.NewLogger()
	auth := factory.NewPrivateGuard("/htmx/login")

	router.HandleFunc("/htmx/todos/remove", log(auth(module.removeTodo)))
	router.HandleFunc("/htmx/todos/toggle", log(auth(module.toggleTodo)))
	router.HandleFunc("/htmx/todos/add", log(auth(module.addTodo)))
	router.HandleFunc("/htmx/todos/list", log(auth(module.getTodoList)))
	router.HandleFunc("/htmx/todos", log(auth(module.getTodoPage)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusMovedPermanently)
	})
}
