package htmx

import "net/http"

func (client *Client) getLogoutPage(response http.ResponseWriter, request *http.Request) {
	loggedIn := client.app.IsLoggedIn(request)
	html := client.renderer.renderLogoutPage(loggedIn)
	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func (client *Client) logoutUser(response http.ResponseWriter, request *http.Request) {
	err := client.app.LogoutUser(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/htmx/logout")
	response.WriteHeader(http.StatusOK)
}
