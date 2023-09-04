package infrastructure

import (
	"net/http"
)

type RouteHandlerFunc func(http.ResponseWriter, *http.Request) error

type RouteHandler interface {
	HandleRoute(http.ResponseWriter, *http.Request) error
}
