package htmx

import (
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
)

type registerPageData struct {
	Key      int64
	Name     string
	Password string
	Error    string
}

type registerPageController struct {
	securityService *security.SecurityService
	*defaultRenderer
}

func newRegisterPageController(security *security.SecurityService) *registerPageController {
	registerPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/register_page.html"))
	return &registerPageController{security, newDefaultRenderer(registerPage)}
}

func (this *registerPageController) page(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.securityService.IsLoggedIn(request)
	if isLoggedIn {
		http.Redirect(response, request, "/htmx/todos", http.StatusSeeOther)
		return
	}

	this.render(response, "page", registerPageData{
		Key: newRenderKey(),
	}, nil)
}

func (this *registerPageController) registerUser(response http.ResponseWriter, request *http.Request) {
	isLoggedIn := this.securityService.IsLoggedIn(request)
	if isLoggedIn {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	name := request.FormValue("name")
	password := request.FormValue("password")
	renderError := func(errorMessage string) {
		this.render(response, "form", registerPageData{
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

	err := this.securityService.RegisterUser(name, password, response)
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
