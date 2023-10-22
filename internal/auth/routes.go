package auth

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

	router.HandleFunc("/htmx/users/logout", log(auth(service.logoutUser)))
	router.HandleFunc("/htmx/users/login", log(service.loginUser))
	router.HandleFunc("/htmx/users/register", log(service.registerUser))
	router.HandleFunc("/htmx/logout", log(service.getLogoutPage))
	router.HandleFunc("/htmx/login", log(service.getLoginPage))
	router.HandleFunc("/htmx/register", log(service.getRegisterPage))
}
