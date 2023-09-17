package middleware

import (
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/users/domain"
)

type SessionStore interface {
	Remove(sessionId string)
	Validate(sessionId string) (*domain.Session, bool)
	Extend(*domain.Session) *domain.Session
}

type RequireSessionMiddleware struct {
	next     http.Handler
	sessions SessionStore
	mode     string
}

func RequireSession(handler http.Handler, sessions SessionStore, mode string) *RequireSessionMiddleware {
	return &RequireSessionMiddleware{handler, sessions, mode}
}

func (middleware *RequireSessionMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	redirectToLogin := func() {
		// page redirect
		if request.Method == http.MethodGet {
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return
		}

		// htmx redirect
		response.Header().Add("HX-Redirect", "/login")
		response.WriteHeader(http.StatusOK)
		return
	}

	sessionCookie, err := request.Cookie(domain.SessionCookieName)
	if err != nil {
		redirectToLogin()
		return
	}

	sessionId := sessionCookie.Value
	session, ok := middleware.sessions.Validate(sessionId)
	if !ok {
		redirectToLogin()
		return
	}

	session = middleware.sessions.Extend(session)
	http.SetCookie(response, domain.NewSessionCookie(session, middleware.mode))
	middleware.next.ServeHTTP(response, request)
}
