package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

func NewClient(securityService *security.SecurityService, todoService *todo.TodoService, router *http.ServeMux) {
	log := newRequestLogger()
	private := newSessionGuard(securityService, "/htmx/login")
	router.Handle(assetPath, newAssetHandler())

	registerPageController := newRegisterPageController(securityService)
	router.HandleFunc("GET /htmx/register", log(registerPageController.page))
	router.HandleFunc("POST /htmx/api/register", log(registerPageController.registerUser))

	loginPageController := newLoginPageController(securityService)
	router.HandleFunc("GET /htmx/login", log(loginPageController.page))
	router.HandleFunc("POST /htmx/api/login", log(loginPageController.loginUser))

	logoutPageController := newLogoutPageController(securityService)
	router.HandleFunc("GET /htmx/logout", log(logoutPageController.page))
	router.HandleFunc("DELETE /htmx/api/logout", log(private(logoutPageController.logoutUser)))

	todoListPageController := newTodoListPageController(todoService)
	router.HandleFunc("GET /htmx/todo-lists", log(private(todoListPageController.page)))
	router.HandleFunc("GET /htmx/api/todo-lists/list", log(private(todoListPageController.lists)))
	router.HandleFunc("POST /htmx/api/todo-lists/add", log(private(todoListPageController.addList)))
	router.HandleFunc("DELETE /htmx/api/todo-lists/remove", log(private(todoListPageController.removeList)))

	todoPageController := newTodoPageController(todoService)
	router.HandleFunc("GET /htmx/todos", log(private(todoPageController.page)))
	router.HandleFunc("GET /htmx/api/todos/list", log(private(todoPageController.todos)))
	router.HandleFunc("POST /htmx/api/todos/add", log(private(todoPageController.addTodo)))
	router.HandleFunc("PATCH /htmx/api/todos/toggle", log(private(todoPageController.toggleTodo)))
	router.HandleFunc("DELETE /htmx/api/todos/remove", log(private(todoPageController.removeTodo)))

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todo-lists", http.StatusSeeOther)
	})
}
