package common

import (
	"log"
	"net/http"
)

type RouteHandler func(http.ResponseWriter, *http.Request) error

func (this RouteHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if err := this(response, request); err != nil {
		switch event := err.(type) {
		case *RouteError:
			log.Println(event.Message)
			http.Error(response, event.Message, event.Code)
		default:
			log.Printf("Internal server error: %v", event)
			http.Error(response, event.Error(), http.StatusInternalServerError)
		}

	}
}
