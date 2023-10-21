package auth

import (
	"net/http"
)

type htmxController struct {
	authenticationService *authenticationService
	htmxRenderer          *htmxRenderer
}

func NewHtmxController(service *authenticationService, renderer *htmxRenderer) *htmxController {
	return &htmxController{service, renderer}
}

func (controller *htmxController) getRegisterPage(response http.ResponseWriter, request *http.Request) {
	html := controller.htmxRenderer.renderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) registerUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := controller.htmxRenderer.renderRegisterForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
		return
	}

	err := controller.authenticationService.registerUser(name, password, response)
	if err != nil {
		renderError("Application error.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (controller *htmxController) getLoginPage(response http.ResponseWriter, request *http.Request) {
	html := controller.htmxRenderer.renderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	html := controller.htmxRenderer.renderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (controller *htmxController) loginUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := controller.htmxRenderer.renderLoginForm(name, password, message)
		response.Header().Add("Content-type", "text/html; charset=utf-8")
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
		return
	}

	err := controller.authenticationService.loginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (controller *htmxController) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := controller.authenticationService.logoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
