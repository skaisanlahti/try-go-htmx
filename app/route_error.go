package app

import (
	"log"
	"net/http"
)

type RouteError struct {
	Message string
	Code    int
}

func (this *RouteError) Error() string {
	return this.Message
}

func (this *RouteError) Log(err error) *RouteError {
	log.Printf("Error: %v", err)
	return this
}

func NewRouteError(msg string, code int) *RouteError {
	return &RouteError{
		Message: msg,
		Code:    code,
	}
}

type ErrorHandlerFunc struct {
	handlerFunc RouteHandlerFunc
}

func NewErrorHandlerFunc(handler RouteHandlerFunc) *ErrorHandlerFunc {
	return &ErrorHandlerFunc{handlerFunc: handler}
}

func (this *ErrorHandlerFunc) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if err := this.handlerFunc(response, request); err != nil {
		switch event := err.(type) {
		case *RouteError:
			log.Println(event.Message)
			http.Error(response, event.Message, event.Code)
		case error:
			log.Printf("Error occurred: %s", event.Error())
		default:
			log.Printf("Internal server error: %v", event)
		}

	}
}

type ErrorHandler struct {
	route RouteHandler
}

func NewErrorHandler(handler RouteHandler) *ErrorHandler {
	return &ErrorHandler{route: handler}
}

func (this *ErrorHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if err := this.route.HandleRoute(response, request); err != nil {
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
