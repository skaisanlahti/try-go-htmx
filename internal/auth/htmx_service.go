package auth

import (
	"net/http"
)

func NewHtmxService(service *authenticationService, renderer *htmxRenderer) *htmxService {
	return &htmxService{service, renderer}
}

type htmxService struct {
	authenticationService *authenticationService
	htmxRenderer          *htmxRenderer
}

func (service *htmxService) getRegisterPage(response http.ResponseWriter, request *http.Request) {
	html := service.htmxRenderer.renderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) registerUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := service.htmxRenderer.renderRegisterForm(name, password, message)
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

	err := service.authenticationService.registerUser(name, password, response)
	if err != nil {
		renderError("Application error.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (service *htmxService) getLoginPage(response http.ResponseWriter, request *http.Request) {
	html := service.htmxRenderer.renderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	html := service.htmxRenderer.renderLogoutPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) loginUser(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := service.htmxRenderer.renderLoginForm(name, password, message)
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

	err := service.authenticationService.loginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (service *htmxService) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := service.authenticationService.logoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
