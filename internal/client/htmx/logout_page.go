package htmx

import (
	"html/template"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/security"
)

type logoutPageData struct {
	LoggedIn bool
}

type logoutPageController struct {
	securityService *security.SecurityService
	*defaultRenderer
}

func newLogoutPageController(security *security.SecurityService) *logoutPageController {
	logoutPage := template.Must(template.ParseFS(templateFiles, "web/html/page.html", "web/html/logout_page.html"))
	return &logoutPageController{security, newDefaultRenderer(logoutPage)}
}

func (this *logoutPageController) page(response http.ResponseWriter, request *http.Request) {
	this.render(response, "page", logoutPageData{
		LoggedIn: this.securityService.IsLoggedIn(request),
	}, nil)
}

func (this *logoutPageController) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := this.securityService.LogoutUser(response, request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
