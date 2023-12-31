package security

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
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
		session: newSessionService(newSessionStorage(), sessionOptions),
		user:    newUserService(newUserStorage(database), newPasswordHasher(passwordOptions)),
	}
}

func (this *SecurityService) RegisterUser(name string, password string, response http.ResponseWriter) error {
	userId, err := this.user.addUser(name, password)
	if err != nil {
		return err
	}

	err = this.session.startSession(response, userId)
	if err != nil {
		return err
	}

	return nil
}

func (this *SecurityService) LoginUser(name string, password string, response http.ResponseWriter) error {
	user, err := this.user.verifyUser(name, password)
	if err != nil {
		return err
	}

	err = this.session.startSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (this *SecurityService) LogoutUser(response http.ResponseWriter, request *http.Request) error {
	err := this.session.clearSession(response, request)
	if err != nil {
		return err
	}

	return nil
}

func (this *SecurityService) IsLoggedIn(request *http.Request) bool {
	ok := this.session.sessionExists(request)
	if !ok {
		return false
	}

	return true
}

func (this *SecurityService) VerifySession(response http.ResponseWriter, request *http.Request) (entity.User, error) {
	var user entity.User
	session, err := this.session.verifySession(response, request)
	if err != nil {
		return user, err
	}
	user, err = this.user.storage.findUserById(session.UserId)
	if err != nil {
		return user, err
	}

	return user, nil
}
