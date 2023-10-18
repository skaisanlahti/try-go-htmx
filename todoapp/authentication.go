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

func CreateUser(name string, key []byte) (User, error) {
	nameLength := len([]rune(name))
	if nameLength < 3 {
		return User{}, errors.New("User name is too short.")
	}

	if nameLength > 100 {
		return User{}, errors.New("User name is too long.")
	}

	return User{Name: name, Key: key}, nil
}

func createFakeUser(encoder PasswordService) User {
	fakeKey, err := encoder.CreateKey("Fake password to compare")
	if err != nil {
		log.Panicln("Failed to create fake key for login.")
	}

	fakeUser, err := CreateUser("FakeName", fakeKey)
	if err != nil {
		log.Panicln("Failed to create fake user for login.")
	}

	return fakeUser
}

type UserStorage interface {
	FindUsers() []User
	FindUserByName(name string) (User, error)
	AddUser(user User) (int, error)
	UpdateUserKey(user User) error
	RemoveUser(id int) error
}

type AuthenticationService struct {
	SessionService  *SessionService
	PasswordService PasswordService
	UserStorage     UserStorage
	fakeUser        User
}

func CreateAuthenticationService(sm *SessionService, pe PasswordService, ua UserStorage) *AuthenticationService {
	return &AuthenticationService{sm, pe, ua, createFakeUser(pe)}
}

func (service *AuthenticationService) RegisterUser(name string, password string, response http.ResponseWriter) error {
	user, err := service.UserStorage.FindUserByName(name)
	if err == nil {
		return err
	}

	key, err := service.PasswordService.CreateKey(password)
	if err != nil {
		return err
	}

	user, err = CreateUser(name, key)
	if err != nil {
		return err
	}

	userId, err := service.UserStorage.AddUser(user)
	if err != nil {
		return err
	}

	err = service.SessionService.StartSession(response, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthenticationService) LoginUser(name string, password string, response http.ResponseWriter) error {
	user, err := service.UserStorage.FindUserByName(name)
	if err != nil {
		service.PasswordService.VerifyKey(service.fakeUser.Key, password)
		return err
	}

	isPasswordCorrect, newKeyChannel := service.PasswordService.VerifyKey(user.Key, password)
	if !isPasswordCorrect {
		return err
	}

	if newKeyChannel != nil {
		go service.updateKey(user, newKeyChannel)
	}

	err = service.SessionService.StartSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthenticationService) updateKey(user User, newKeyChannel <-chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := service.UserStorage.UpdateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}

func (service *AuthenticationService) LogoutUser(response http.ResponseWriter, request *http.Request) error {
	err := service.SessionService.StopSession(response, request)
	if err != nil {
		return err
	}

	return nil
}
