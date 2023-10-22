package auth

import (
	"net/http"
)

type htmxService struct {
	authenticationService *authenticationService
	htmxRenderer          *htmxRenderer
}

func newHtmxService(service *authenticationService, renderer *htmxRenderer) *htmxService {
	return &htmxService{service, renderer}
}

func (service *htmxService) getRegisterPage(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := service.authenticationService.sessionService.sessionExists(request)
	if isLoggedIn {
		http.Redirect(response, request, "htmx/todos", http.StatusSeeOther)
		return
	}

	html := service.htmxRenderer.renderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) registerUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := service.authenticationService.sessionService.sessionExists(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

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
	if err == ErrUserAlreadyExists {
		renderError(err.Error())
		return
	}

	if err != nil {
		renderError("Something went wrong.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}

func (service *htmxService) getLoginPage(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := service.authenticationService.sessionService.sessionExists(request)
	if isLoggedIn {
		http.Redirect(response, request, "htmx/todos", http.StatusSeeOther)
		return
	}

	html := service.htmxRenderer.renderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (service *htmxService) loginUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := service.authenticationService.sessionService.sessionExists(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (service *htmxService) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	loggedIn := service.authenticationService.sessionService.sessionExists(request)
	html := service.htmxRenderer.renderLogoutPage(loggedIn)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
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
