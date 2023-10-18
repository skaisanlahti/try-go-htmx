package htmx

import (
	"net/http"
)

func mapRoutes(router *http.ServeMux, client *Client) {
	requireSession := createSessionGuard(client.AuthenticationService.SessionService, "/htmx/login")
	logRequest := createRequestLogger()

	addAssets(router)

	router.HandleFunc("/htmx/users/logout", logRequest(requireSession(client.logoutUser)))
	router.HandleFunc("/htmx/users/login", logRequest(client.loginUser))
	router.HandleFunc("/htmx/users/register", logRequest(client.registerUser))
	router.HandleFunc("/htmx/logout", logRequest(client.getLogoutPage))
	router.HandleFunc("/htmx/login", logRequest(client.getLoginPage))
	router.HandleFunc("/htmx/register", logRequest(client.getRegisterPage))

	router.HandleFunc("/htmx/todos/remove", logRequest(requireSession(client.removeTodo)))
	router.HandleFunc("/htmx/todos/toggle", logRequest(requireSession(client.toggleTodo)))
	router.HandleFunc("/htmx/todos/add", logRequest(requireSession(client.addTodo)))
	router.HandleFunc("/htmx/todos/list", logRequest(requireSession(client.getTodoList)))
	router.HandleFunc("/htmx/todos", logRequest(requireSession(client.getTodoPage)))

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/htmx/todos", http.StatusMovedPermanently)
	})
}
