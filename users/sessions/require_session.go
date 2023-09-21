package sessions

import (
	"net/http"
)

type SessionVerifier interface {
	VerifySession(response http.ResponseWriter, request *http.Request) error
}

type RequireSessionMiddleware struct {
	next    http.Handler
	session SessionVerifier
}

func RequireSession(handler http.Handler, sessions SessionVerifier) *RequireSessionMiddleware {
	return &RequireSessionMiddleware{handler, sessions}
}

func (middleware *RequireSessionMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := middleware.session.VerifySession(response, request)
	if err != nil {
		// page redirect
		if request.Method == http.MethodGet {
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return
		}

		// htmx redirect
		response.Header().Add("HX-Location", "/login")
		response.WriteHeader(http.StatusOK)
		return
	}

	middleware.next.ServeHTTP(response, request)
}
