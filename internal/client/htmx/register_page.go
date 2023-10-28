package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
)

func (this *Client) getRegisterPage(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.app.IsLoggedIn(request)
	if isLoggedIn {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
		return
	}

	html := this.renderer.renderRegisterPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) registerUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.app.IsLoggedIn(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := this.renderer.renderRegisterForm(name, password, message)
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

	err := this.app.RegisterUser(name, password, response)
	if err == security.ErrUserAlreadyExists {
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
