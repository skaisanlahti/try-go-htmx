package htmx

import (
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (proxy *responseRecorder) WriteHeader(code int) {
	proxy.statusCode = code
	proxy.ResponseWriter.WriteHeader(code)
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

func newSessionGuard(session *todoapp.SessionManager, redirectUrl string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			err := session.VerifySession(response, request)
			if err != nil {
				// page redirect
				if request.Method == http.MethodGet {
					http.Redirect(response, request, redirectUrl, http.StatusSeeOther)
					return
				}

				// htmx redirect
				response.Header().Add("HX-Location", redirectUrl)
				response.WriteHeader(http.StatusOK)
				return
			}

			next(response, request)
		}
	}
}
