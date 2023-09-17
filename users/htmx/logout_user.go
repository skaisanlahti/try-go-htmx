package htmx

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type LogoutSessionStore interface {
	Remove(sessionId string)
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
	cookie, err := request.Cookie("sid")
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId := cookie.Value
	handler.sessions.Remove(sessionId)
	http.SetCookie(response, domain.NewExpiredSessionCookie())
	response.Header().Add("HX-Redirect", "/logout")
	response.WriteHeader(http.StatusOK)
}
