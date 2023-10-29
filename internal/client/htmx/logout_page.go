package htmx

import "net/http"

func (this *controller) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	loggedIn := this.model.IsLoggedIn(request)
	html := this.view.renderLogoutPage(loggedIn)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *controller) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := this.model.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
