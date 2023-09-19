package handlers

import (
	"net/http"
)

type LogoutSessionStore interface {
	Remove(request *http.Request) (*http.Cookie, error)
}

type LogoutUserHandler struct {
	sessions LogoutSessionStore
}

func NewLogoutUserHandler(
	sessions LogoutSessionStore,
) *LogoutUserHandler {
	return &LogoutUserHandler{sessions}
}

func (handler *LogoutUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	cookie, err := handler.sessions.Remove(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(response, cookie)
	response.Header().Add("HX-Location", "/logout")
	response.WriteHeader(http.StatusOK)
}
