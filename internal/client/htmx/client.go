package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

func NewClient(security *security.SecurityService, todo *todo.TodoService, router *http.ServeMux) {
	controller := newController(security, todo)
	log := controller.logRequest()
	private := controller.requireSession("/htmx/login")

	router.HandleFunc(assetPath, controller.getAssets)

	router.HandleFunc("/htmx/register", log(controller.getRegisterPage))
	router.HandleFunc("/htmx/api/register", log(controller.registerUser))

	router.HandleFunc("/htmx/login", log(controller.getLoginPage))
	router.HandleFunc("/htmx/api/login", log(controller.loginUser))

	router.HandleFunc("/htmx/logout", log(controller.getLogoutPage))
	router.HandleFunc("/htmx/api/logout", log(private(controller.logoutUser)))

	router.HandleFunc("/htmx/todos", log(private(controller.getTodoPage)))
	router.HandleFunc("/htmx/api/todos/list", log(private(controller.getTodoList)))
	router.HandleFunc("/htmx/api/todos/add", log(private(controller.addTodo)))
	router.HandleFunc("/htmx/api/todos/toggle", log(private(controller.toggleTodo)))
	router.HandleFunc("/htmx/api/todos/remove", log(private(controller.removeTodo)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
	})
}
