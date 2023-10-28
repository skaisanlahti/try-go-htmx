package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

type Client struct {
	app      *application
	renderer *renderer
}

func NewClient(security *security.SecurityService, todo *todo.TodoService, router *http.ServeMux) {
	client := &Client{&application{security, todo}, newRenderer()}
	log := client.logRequest()
	private := client.requireSession("/htmx/login")

	useAssets(router)

	router.HandleFunc("/htmx/register", log(client.getRegisterPage))
	router.HandleFunc("/htmx/api/register", log(client.registerUser))

	router.HandleFunc("/htmx/login", log(client.getLoginPage))
	router.HandleFunc("/htmx/api/login", log(client.loginUser))

	router.HandleFunc("/htmx/logout", log(client.getLogoutPage))
	router.HandleFunc("/htmx/api/logout", log(private(client.logoutUser)))

	router.HandleFunc("/htmx/todos", log(private(client.getTodoPage)))
	router.HandleFunc("/htmx/api/todos/list", log(private(client.getTodoList)))
	router.HandleFunc("/htmx/api/todos/add", log(private(client.addTodo)))
	router.HandleFunc("/htmx/api/todos/toggle", log(private(client.toggleTodo)))
	router.HandleFunc("/htmx/api/todos/remove", log(private(client.removeTodo)))
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
	})
}
