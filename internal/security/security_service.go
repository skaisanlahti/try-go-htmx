package security

import (
	"database/sql"
	"net/http"
)

type SecurityService struct {
	session *sessionService
	user    *userService
}

func NewSecurityService(
	database *sql.DB,
	passwordOptions PasswordOptions,
	sessionOptions SessionOptions,
) *SecurityService {
	return &SecurityService{
		session: newSessionService(sessionOptions, newSessionStorage()),
		user:    newUserService(newUserStorage(database), newPasswordHasher(passwordOptions)),
	}
}

func (service *SecurityService) RegisterUser(name string, password string, response http.ResponseWriter) error {
	userId, err := service.user.newUser(name, password)
	if err != nil {
		return err
	}

	err = service.session.startSession(response, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *SecurityService) LoginUser(name string, password string, response http.ResponseWriter) error {
	user, err := service.user.verifyUser(name, password)
	if err != nil {
		return err
	}

	err = service.session.startSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (service *SecurityService) LogoutUser(response http.ResponseWriter, request *http.Request) error {
	err := service.session.clearSession(response, request)
	if err != nil {
		return err
	}

	return nil
}

func (service *SecurityService) IsLoggedIn(request *http.Request) bool {
	ok := service.session.sessionExists(request)
	if !ok {
		return false
	}

	return true
}

func (service *SecurityService) VerifySession(response http.ResponseWriter, request *http.Request) error {
	err := service.session.verifySession(response, request)
	if err != nil {
		return err
	}

	return nil
}
