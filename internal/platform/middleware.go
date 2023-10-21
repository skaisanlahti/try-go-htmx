package platform

import (
	"log"
	"net/http"
	"time"
)

type sessionVerifier interface {
	VerifySession(http.ResponseWriter, *http.Request) error
}

type middlewareFactory struct {
	sessionVerifier sessionVerifier
}

func NewMiddlewareFactory(sessionVerifier sessionVerifier) *middlewareFactory {
	return &middlewareFactory{sessionVerifier}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (recorder *responseRecorder) WriteHeader(code int) {
	recorder.statusCode = code
	recorder.ResponseWriter.WriteHeader(code)
}

func (middleware *middlewareFactory) NewLogger() func(http.HandlerFunc) http.HandlerFunc {
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

func (middleware *middlewareFactory) NewSessionGuard(redirectUrl string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			err := middleware.sessionVerifier.VerifySession(response, request)
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
