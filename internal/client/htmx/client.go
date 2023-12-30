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
	router.HandleFunc("/htmx/register", log(registerPageController.page))
	router.HandleFunc("/htmx/api/register", log(registerPageController.registerUser))

	loginPageController := newLoginPageController(securityService)
	router.HandleFunc("/htmx/login", log(loginPageController.page))
	router.HandleFunc("/htmx/api/login", log(loginPageController.loginUser))

	logoutPageController := newLogoutPageController(securityService)
	router.HandleFunc("/htmx/logout", log(logoutPageController.page))
	router.HandleFunc("/htmx/api/logout", log(private(logoutPageController.logoutUser)))

	todoListPageController := newTodoListPageController(todoService)
	router.HandleFunc("/htmx/todo-lists", log(private(todoListPageController.page)))
	router.HandleFunc("/htmx/api/todo-lists/list", log(private(todoListPageController.lists)))
	router.HandleFunc("/htmx/api/todo-lists/add", log(private(todoListPageController.addList)))
	router.HandleFunc("/htmx/api/todo-lists/remove", log(private(todoListPageController.removeList)))

	todoPageController := newTodoPageController(todoService)
	router.HandleFunc("/htmx/todos", log(private(todoPageController.page)))
	router.HandleFunc("/htmx/api/todos/list", log(private(todoPageController.todos)))
	router.HandleFunc("/htmx/api/todos/add", log(private(todoPageController.addTodo)))
	router.HandleFunc("/htmx/api/todos/toggle", log(private(todoPageController.toggleTodo)))
	router.HandleFunc("/htmx/api/todos/remove", log(private(todoPageController.removeTodo)))

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todo-lists", http.StatusSeeOther)
	})
}
