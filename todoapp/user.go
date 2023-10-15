package todoapp

import (
	"errors"
	"log"
	"net/http"
)

type User struct {
	Id   int
	Name string
	Key  []byte
}

var (
	ErrUserNameTooShort = errors.New("User name is too short.")
	ErrUserNameTooLong  = errors.New("User name is too long.")
)

func NewUser(name string, key []byte) (User, error) {
	nameLength := len([]rune(name))
	if nameLength < 3 {
		return User{}, ErrUserNameTooShort
	}

	if nameLength > 100 {
		return User{}, ErrUserNameTooLong
	}

	return User{Name: name, Key: key}, nil
}

type UserAccessor interface {
	FindUsers() []User
	FindUserByName(name string) (User, error)
	AddUser(user User) (int, error)
	UpdateUserKey(user User) error
	RemoveUser(id int) error
}

type UserAuthenticator struct {
	Settings        Settings
	SessionManager  *SessionManager
	PasswordEncoder PasswordEncoder
	UserAccessor    UserAccessor
	fakeUser        User
}

func newFakeUser(encoder PasswordEncoder) User {
	fakeKey, err := encoder.NewKey("Fake password to compare")
	if err != nil {
		log.Panicln("Failed to create fake key for login.")
	}

	fakeUser, err := NewUser("FakeName", fakeKey)
	if err != nil {
		log.Panicln("Failed to create fake user for login.")
	}

	return fakeUser
}

func NewUserAuthenticator(s Settings, sm *SessionManager, pe PasswordEncoder, ua UserAccessor) *UserAuthenticator {
	return &UserAuthenticator{s, sm, pe, ua, newFakeUser(pe)}
}

func (authenticator *UserAuthenticator) RegisterUser(name string, password string, response http.ResponseWriter) error {
	newUser, err := authenticator.UserAccessor.FindUserByName(name)
	if err == nil {
		return err
	}

	key, err := authenticator.PasswordEncoder.NewKey(password)
	if err != nil {
		return err
	}

	newUser, err = NewUser(name, key)
	if err != nil {
		return err
	}

	userId, err := authenticator.UserAccessor.AddUser(newUser)
	if err != nil {
		return err
	}

	err = authenticator.SessionManager.StartSession(response, userId)
	if err != nil {
		return err
	}

	return nil
}

func (authenticator *UserAuthenticator) LoginUser(name string, password string, response http.ResponseWriter) error {
	user, err := authenticator.UserAccessor.FindUserByName(name)
	if err != nil {
		authenticator.PasswordEncoder.VerifyKey(authenticator.fakeUser.Key, password)
		return err
	}

	isPasswordCorrect, newKeyChannel := authenticator.PasswordEncoder.VerifyKey(user.Key, password)
	if !isPasswordCorrect {
		return err
	}

	if newKeyChannel != nil {
		go authenticator.updateKey(user, newKeyChannel)
	}

	err = authenticator.SessionManager.StartSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (authenticator *UserAuthenticator) updateKey(user User, newKeyChannel <-chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := authenticator.UserAccessor.UpdateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}

func (authenticator *UserAuthenticator) LogoutUser(response http.ResponseWriter, request *http.Request) error {
	err := authenticator.SessionManager.StopSession(response, request)
	if err != nil {
		return err
	}

	return nil
}
