package htmx

import (
	"net/http"
)

func AddRoutes(router *http.ServeMux, client *Client) {
	requireSession := newSessionGuard(client.UserAuthenticator.SessionManager, "/htmx/login")
	logRequest := newRequestLogger()

	AddAssets(router)

	router.HandleFunc("/htmx/users/logout", logRequest(requireSession(client.LogoutUser)))
	router.HandleFunc("/htmx/users/login", logRequest(client.LoginUser))
	router.HandleFunc("/htmx/users/register", logRequest(client.RegisterUser))
	router.HandleFunc("/htmx/logout", logRequest(client.GetLogoutPage))
	router.HandleFunc("/htmx/login", logRequest(client.GetLoginPage))
	router.HandleFunc("/htmx/register", logRequest(client.GetRegisterPage))

	router.HandleFunc("/htmx/todos/remove", logRequest(requireSession(client.RemoveTodo)))
	router.HandleFunc("/htmx/todos/toggle", logRequest(requireSession(client.ToggleTodo)))
	router.HandleFunc("/htmx/todos/add", logRequest(requireSession(client.AddTodo)))
	router.HandleFunc("/htmx/todos/list", logRequest(requireSession(client.GetTodoList)))
	router.HandleFunc("/htmx/todos", logRequest(requireSession(client.GetTodoPage)))

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusMovedPermanently)
	})
}
