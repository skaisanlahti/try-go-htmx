package auth

import (
	"log"
	"net/http"
)

type authenticationService struct {
	sessionService  *sessionService
	passwordService *passwordHasher
	userStorage     *userStorage
	fakeUser        user
}

func newAuthenticationService(
	sessionService *sessionService,
	passwordService *passwordHasher,
	userStorage *userStorage,
) *authenticationService {
	fakeKey, err := passwordService.hash("password")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fakeUser, _ := newUser("username", fakeKey)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &authenticationService{sessionService, passwordService, userStorage, fakeUser}
}

func (service *authenticationService) registerUser(name string, password string, response http.ResponseWriter) error {
	user, err := service.userStorage.findUserByName(name)
	if err == nil {
		return err
	}

	key, err := service.passwordService.hash(password)
	if err != nil {
		return err
	}

	user, err = newUser(name, key)
	if err != nil {
		return err
	}

	userId, err := service.userStorage.addUser(user)
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
	user, err := service.userStorage.findUserByName(name)
	if err != nil {
		service.passwordService.verify(service.fakeUser.Key, password)
		return err
	}

	isPasswordCorrect, newKeyChannel := service.passwordService.verify(user.Key, password)
	if !isPasswordCorrect {
		return err
	}

	if newKeyChannel != nil {
		go service.updateKey(user, newKeyChannel)
	}

	err = service.sessionService.startSession(response, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (service *authenticationService) updateKey(user user, newKeyChannel <-chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := service.userStorage.updateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}

func (service *authenticationService) logoutUser(response http.ResponseWriter, request *http.Request) error {
	err := service.sessionService.stopSession(response, request)
	if err != nil {
		return err
	}

	return nil
}
