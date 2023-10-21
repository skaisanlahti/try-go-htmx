package todo

import (
	"net/http"
)

type middlewareFactory interface {
	NewLogger() func(http.HandlerFunc) http.HandlerFunc
	NewSessionGuard(redirectUrl string) func(http.HandlerFunc) http.HandlerFunc
}

func MapRoutes(router *http.ServeMux, client *htmxController, middleware middlewareFactory) {
	log := middleware.NewLogger()
	auth := middleware.NewSessionGuard("/htmx/login")

	router.HandleFunc("/htmx/todos/remove", log(auth(client.removeTodo)))
	router.HandleFunc("/htmx/todos/toggle", log(auth(client.toggleTodo)))
	router.HandleFunc("/htmx/todos/add", log(auth(client.addTodo)))
	router.HandleFunc("/htmx/todos/list", log(auth(client.getTodoList)))
	router.HandleFunc("/htmx/todos", log(auth(client.getTodoPage)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusMovedPermanently)
	})
}
