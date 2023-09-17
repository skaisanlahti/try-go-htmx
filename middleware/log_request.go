package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseProxy struct {
	http.ResponseWriter
	statusCode int
}

func (proxy *responseProxy) WriteHeader(code int) {
	proxy.statusCode = code
	proxy.ResponseWriter.WriteHeader(code)
}

type LogRequestMiddleware struct {
	next http.Handler
}

func LogRequest(handler http.Handler) *LogRequestMiddleware {
	return &LogRequestMiddleware{next: handler}
}

func (middleware *LogRequestMiddleware) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	before := time.Now().UnixMilli()
	responseProxy := &responseProxy{ResponseWriter: response}

	middleware.next.ServeHTTP(responseProxy, request)

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
