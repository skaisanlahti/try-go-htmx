package htmx

import "net/http"

func (this *Client) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	loggedIn := this.app.IsLoggedIn(request)
	html := this.renderer.renderLogoutPage(loggedIn)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (this *Client) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := this.app.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
