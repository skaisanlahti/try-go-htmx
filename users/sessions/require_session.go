package sessions

import (
	"net/http"
)

type SessionStore interface {
	Validate(request *http.Request) (*Session, error)
	Extend(*Session) (*http.Cookie, error)
}

type RequireSessionMiddleware struct {
	next     http.Handler
	sessions SessionStore
}

func RequireSession(handler http.Handler, sessions SessionStore) *RequireSessionMiddleware {
	return &RequireSessionMiddleware{handler, sessions}
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

	session, err := middleware.sessions.Validate(request)
	if err != nil {
		redirectToLogin()
		return
	}

	newSessionCookie, err := middleware.sessions.Extend(session)
	if err != nil {
		redirectToLogin()
		return
	}

	http.SetCookie(response, newSessionCookie)
	middleware.next.ServeHTTP(response, request)
}
