package auth

import (
	"net/http"
)

type authenticationService struct {
	sessionService *sessionService
	userService    *userService
}

func newAuthenticationService(
	sessionService *sessionService,
	userService *userService,
) *authenticationService {
	return &authenticationService{sessionService, userService}
}

func (service *authenticationService) registerUser(name string, password string, response http.ResponseWriter) error {
	userId, err := service.userService.newUser(name, password)
	if err != nil {
		return err
	}

	err = service.sessionService.startSession(response, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *authenticationService) loginUser(name string, password string, response http.ResponseWriter) error {
	user, err := service.userService.verifyUser(name, password)
	if err != nil {
		return err
	}

	err = service.sessionService.startSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (service *authenticationService) logoutUser(response http.ResponseWriter, request *http.Request) error {
	err := service.sessionService.stopSession(response, request)
	if err != nil {
		return err
	}

	return nil
}
