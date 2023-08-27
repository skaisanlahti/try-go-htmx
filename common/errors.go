package common

import "log"

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
