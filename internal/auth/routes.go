package auth

import (
	"net/http"
)

type middlewareFactory interface {
	NewLogger() func(http.HandlerFunc) http.HandlerFunc
	NewSessionGuard(redirectUrl string) func(http.HandlerFunc) http.HandlerFunc
}

func MapRoutes(router *http.ServeMux, controller *htmxController, middleware middlewareFactory) {
	log := middleware.NewLogger()
	auth := middleware.NewSessionGuard("/htmx/login")

	router.HandleFunc("/htmx/users/logout", log(auth(controller.logoutUser)))
	router.HandleFunc("/htmx/users/login", log(controller.loginUser))
	router.HandleFunc("/htmx/users/register", log(controller.registerUser))
	router.HandleFunc("/htmx/logout", log(controller.getLogoutPage))
	router.HandleFunc("/htmx/login", log(controller.getLoginPage))
	router.HandleFunc("/htmx/register", log(controller.getRegisterPage))
}
