package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (proxy *statusWriter) WriteHeader(code int) {
	proxy.statusCode = code
	proxy.ResponseWriter.WriteHeader(code)
}

type Logger struct {
	next http.Handler
}

func Log(handler http.Handler) *Logger {
	return &Logger{next: handler}
}

func (handler *Logger) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	before := time.Now().UnixMilli()
	responseProxy := &statusWriter{ResponseWriter: response}

	handler.next.ServeHTTP(responseProxy, request)

	after := time.Now().UnixMilli()
	duration := after - before
	log.Printf(
		"%s %s | status %d | duration %d ms",
		request.Method,
		request.URL.Path,
		responseProxy.statusCode,
		duration,
	)
}
