package handlers

import (
	"net/http"
)

type LogoutUserSession interface {
	ClearSession(response http.ResponseWriter, request *http.Request) error
}

type LogoutUserHandler struct {
	session LogoutUserSession
}

func NewLogoutUserHandler(
	session LogoutUserSession,
) *LogoutUserHandler {
	return &LogoutUserHandler{session}
}

func (handler *LogoutUserHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := handler.session.ClearSession(response, request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Add("HX-Location", "/logout")
	response.WriteHeader(http.StatusOK)
}
