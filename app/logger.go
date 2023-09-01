package app

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (this *statusWriter) WriteHeader(code int) {
	this.statusCode = code
	this.ResponseWriter.WriteHeader(code)
}

type Logger struct {
	next http.Handler
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{next: handler}
}

func (this *Logger) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	before := time.Now().UnixMilli()
	responseWithStatus := &statusWriter{ResponseWriter: response}

	this.next.ServeHTTP(responseWithStatus, request)

	after := time.Now().UnixMilli()
	duration := after - before
	log.Printf(
		"%s %s | status %d | duration %d ms",
		request.Method,
		request.URL.Path,
		responseWithStatus.statusCode,
		duration,
	)
}
