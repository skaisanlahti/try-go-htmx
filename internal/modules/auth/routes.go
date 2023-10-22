package auth

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/platform"
)

func MapRoutes(router *http.ServeMux, module *authModule, factory platform.MiddlewareFactory) {
	log := factory.NewLogger()
	auth := factory.NewPrivateGuard("/htmx/login")

	router.HandleFunc("/htmx/users/logout", log(auth(module.logoutUser)))
	router.HandleFunc("/htmx/logout", log(module.getLogoutPage))
	router.HandleFunc("/htmx/users/login", log(module.loginUser))
	router.HandleFunc("/htmx/users/register", log(module.registerUser))
	router.HandleFunc("/htmx/login", log(module.getLoginPage))
	router.HandleFunc("/htmx/register", log(module.getRegisterPage))
}
