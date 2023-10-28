package htmx

import "net/http"

func (this *Client) getLoginPage(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.app.IsLoggedIn(request)
	if isLoggedIn {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
		return
	}

	html := this.renderer.renderLoginPage()
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) loginUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.app.IsLoggedIn(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(message string) {
		html := this.renderer.renderLoginForm(name, password, message)
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

	err := this.app.LoginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}
