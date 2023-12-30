package htmx

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
	"github.com/skaisanlahti/try-go-htmx/internal/security"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (recorder *responseRecorder) WriteHeader(code int) {
	recorder.statusCode = code
	recorder.ResponseWriter.WriteHeader(code)
}

func newRequestLogger() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			start := time.Now()
			recorder := &responseRecorder{ResponseWriter: response}
			next(recorder, request)
			duration := time.Since(start).Milliseconds()
			log.Printf(
				"%s %s | status %d | duration %d ms",
				request.Method,
				request.URL.Path,
				recorder.statusCode,
				duration,
			)
		}
	}
}

func newSessionGuard(security *security.SecurityService, redirectUrl string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			user, err := security.VerifySession(response, request)
			if err != nil {
				if request.Header.Get("HX-Request") == "true" {
					response.Header().Add("HX-Location", redirectUrl)
					response.WriteHeader(http.StatusOK)
					return
				}

				http.Redirect(response, request, redirectUrl, http.StatusSeeOther)
				return
			}

			requestWithUser := addUserToContext(user, request)
			next(response, requestWithUser)
		}
	}
}

func addUserToContext(user entity.User, request *http.Request) *http.Request {
	return request.WithContext(context.WithValue(request.Context(), "user", user))
}

func extractUserFromContext(request *http.Request) (entity.User, bool) {
	user, ok := request.Context().Value("user").(entity.User)
	return user, ok
}
