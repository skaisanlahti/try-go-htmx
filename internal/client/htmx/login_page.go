package htmx

import (
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
)

type loginPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

type loginPageController struct {
	securityService *security.SecurityService
	*defaultRenderer
}

func newLoginPageController(security *security.SecurityService) *loginPageController {
	loginPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/login_page.html"))
	return &loginPageController{security, newDefaultRenderer(loginPage)}
}

func (this *loginPageController) page(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.securityService.IsLoggedIn(request)
	if isLoggedIn {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
		return
	}

	this.render(response, "page", loginPageData{
		Key: newRenderKey(),
	}, nil)
}

func (this *loginPageController) loginUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.securityService.IsLoggedIn(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(errorMessage string) {
		this.render(response, "form", loginPageData{
			Key:      newRenderKey(),
			Name:     name,
			Password: password,
			Error:    errorMessage,
		}, nil)
	}

	if name == "" {
		renderError("Username is required.")
		return
	}

	if password == "" {
		renderError("Password is required.")
		return
	}

	err := this.securityService.LoginUser(name, password, response)
	if err != nil {
		renderError("Invalid credentials.")
		return
	}

	response.Header().Add("HX-Location", "/htmx/todos")
	response.WriteHeader(http.StatusOK)
}
