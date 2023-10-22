package todo

import (
	"net/http"
)

type middlewareFactory interface {
	NewLogger() func(http.HandlerFunc) http.HandlerFunc
	NewSessionGuard(redirectUrl string) func(http.HandlerFunc) http.HandlerFunc
}

func MapRoutes(router *http.ServeMux, service *htmxService, middleware middlewareFactory) {
	log := middleware.NewLogger()
	auth := middleware.NewSessionGuard("/htmx/login")

	router.HandleFunc("/htmx/todos/remove", log(auth(service.removeTodo)))
	router.HandleFunc("/htmx/todos/toggle", log(auth(service.toggleTodo)))
	router.HandleFunc("/htmx/todos/add", log(auth(service.addTodo)))
	router.HandleFunc("/htmx/todos/list", log(auth(service.getTodoList)))
	router.HandleFunc("/htmx/todos", log(auth(service.getTodoPage)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusMovedPermanently)
	})
}
