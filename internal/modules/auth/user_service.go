package auth

import (
	"errors"
	"log"
)

type userService struct {
	userStorage    *userStorage
	passwordHasher *passwordHasher
	fakeUser       user
}

func newUserService(userStorage *userStorage, passwordHasher *passwordHasher) *userService {
	fakeKey, err := passwordHasher.hash("password")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fakeUser, _ := newUser("username", fakeKey)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &userService{userStorage, passwordHasher, fakeUser}
}

var ErrUserAlreadyExists = errors.New("Username already exists.")

func (service *userService) newUser(name string, password string) (int, error) {
	nameExists := service.userStorage.userExists(name)
	if nameExists {
		return 0, ErrUserAlreadyExists
	}

	key, err := service.passwordHasher.hash(password)
	if err != nil {
		return 0, err
	}

	user, err := newUser(name, key)
	if err != nil {
		return 0, err
	}

	userId, err := service.userStorage.addUser(user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

var ErrInvalidCredentials = errors.New("Invalid credentials.")

func (service *userService) verifyUser(name string, password string) (user, error) {
	user, err := service.userStorage.findUserByName(name)
	if err != nil {
		service.passwordHasher.verify(service.fakeUser.Key, password)
		return user, ErrInvalidCredentials
	}

	isPasswordCorrect, newKeyChannel := service.passwordHasher.verify(user.Key, password)
	if !isPasswordCorrect {
		return user, ErrInvalidCredentials
	}

	if newKeyChannel != nil {
		go service.updateUserKey(user, newKeyChannel)
	}

	return user, nil
}

func (service *userService) updateUserKey(user user, newKeyChannel <-chan []byte) {
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
