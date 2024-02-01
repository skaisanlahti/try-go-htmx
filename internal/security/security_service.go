package security

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type SecurityService struct {
	database        *sql.DB
	sessionOptions  SessionOptions
	passwordOptions PasswordOptions
	sessionStorage  *SessionStorage
	cookieFactory   *CookieFactory
	userStorage     *UserStorage
	sessionSigner   *SessionSigner
	passwordHasher  *PasswordHasher
	fakeUser        entity.User
}

func NewSecurityService(
	database *sql.DB,
	passwordOptions PasswordOptions,
	sessionOptions SessionOptions,
) *SecurityService {
	passwordHasher := NewPasswordHasher(passwordOptions)
	fakeKey, err := passwordHasher.Hash("password")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fakeUser := entity.NewUser("username", fakeKey)
	return &SecurityService{
		database:        database,
		sessionOptions:  sessionOptions,
		passwordOptions: passwordOptions,
		cookieFactory:   NewCookieFactory(sessionOptions),
		sessionStorage:  NewSessionStorage(),
		userStorage:     NewUserStorage(database),
		sessionSigner:   NewSessionSigner(sessionOptions),
		passwordHasher:  passwordHasher,
		fakeUser:        fakeUser,
	}
}

func (this *SecurityService) RegisterUser(name string, password string, response http.ResponseWriter) error {
	key, err := this.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	user := entity.NewUser(name, key)
	err = user.Validate()
	if err != nil {
		return err
	}

	userId, err := this.userStorage.InsertUserIfNotExists(user)
	if err != nil {
		return err
	}

	session := entity.NewSession(userId, this.sessionOptions.Duration)
	err = this.sessionStorage.InsertSession(session)
	if err != nil {
		return err
	}

	signedSession, err := this.sessionSigner.NewSignature(session.Id)
	if err != nil {
		return err
	}

	cookie := this.cookieFactory.NewSessionCookie(signedSession)
	http.SetCookie(response, cookie)
	return nil
}

func (this *SecurityService) LoginUser(name string, password string, response http.ResponseWriter) error {
	user, err := this.userStorage.FindUserByName(name)
	if err != nil {
		this.passwordHasher.Verify(this.fakeUser.Key, password)
		return ErrInvalidCredentials
	}

	isPasswordCorrect, newKeyChannel := this.passwordHasher.Verify(user.Key, password)
	if !isPasswordCorrect {
		return ErrInvalidCredentials
	}

	if newKeyChannel != nil {
		go this.updateUserKey(user, newKeyChannel)
	}

	session := entity.NewSession(user.Id, this.sessionOptions.Duration)
	err = this.sessionStorage.InsertSession(session)
	if err != nil {
		return err
	}

	signedSession, err := this.sessionSigner.NewSignature(session.Id)
	if err != nil {
		return err
	}

	cookie := this.cookieFactory.NewSessionCookie(signedSession)
	http.SetCookie(response, cookie)
	return nil
}

func (this *SecurityService) updateUserKey(user entity.User, newKeyChannel chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := this.userStorage.UpdateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}

func (this *SecurityService) LogoutUser(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(this.sessionOptions.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := this.sessionSigner.VerifySignature(cookie.Value)
	if err != nil {
		return err
	}

	err = this.sessionStorage.DeleteSession(sessionId)
	if err != nil {
		return err
	}

	http.SetCookie(response, this.cookieFactory.ClearSessionCookie())
	return nil
}

func (this *SecurityService) IsLoggedIn(request *http.Request) bool {
	cookie, err := request.Cookie(this.sessionOptions.CookieName)
	if err != nil {
		return false
	}

	sessionId, err := this.sessionSigner.VerifySignature(cookie.Value)
	if err != nil {
		return false
	}

	_, err = this.sessionStorage.FindSessionBySessionId(sessionId)
	if err != nil {
		return false
	}

	return true
}

func (this *SecurityService) VerifySession(response http.ResponseWriter, request *http.Request) (entity.User, error) {
	var user entity.User
	cookie, err := request.Cookie(this.sessionOptions.CookieName)
	if err != nil {
		return user, err
	}

	sessionId, err := this.sessionSigner.VerifySignature(cookie.Value)
	if err != nil {
		return user, err
	}

	session, err := this.sessionStorage.FindSessionBySessionId(sessionId)
	if err != nil {
		return user, err
	}

	if session.Expires.Before(time.Now()) {
		this.sessionStorage.DeleteSession(sessionId)
		return user, errors.New("Session has expired.")
	}

	err = this.sessionStorage.UpdateSession(session.Extend(this.sessionOptions.Duration))
	if err != nil {
		return user, err
	}

	signedSession, err := this.sessionSigner.NewSignature(session.Id)
	if err != nil {
		return user, err
	}

	newCookie := this.cookieFactory.NewSessionCookie(signedSession)
	http.SetCookie(response, newCookie)

	user, err = this.userStorage.FindUserById(session.UserId)
	if err != nil {
		return user, err
	}

	return user, nil
}
