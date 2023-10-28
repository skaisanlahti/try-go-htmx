package security

import (
	"errors"
	"log"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type userService struct {
	storage *userStorage
	hasher  *passwordHasher
	fake    entity.User
}

func newUserService(storage *userStorage, password *passwordHasher) *userService {
	fakeKey, err := password.hash("password")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fakeUser := entity.User{Name: "username", Key: fakeKey}
	return &userService{storage, password, fakeUser}
}

var ErrUserAlreadyExists = errors.New("Username already exists.")

func (service *userService) newUser(name string, password string) (int, error) {
	nameExists := service.storage.userExists(name)
	if nameExists {
		return 0, ErrUserAlreadyExists
	}

	key, err := service.hasher.hash(password)
	if err != nil {
		return 0, err
	}

	user := entity.NewUser(name, key)
	err = service.validateUser(user)
	if err != nil {
		return 0, err
	}

	userId, err := service.storage.insertUser(user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (service *userService) validateUser(user entity.User) error {
	nameLength := len([]rune(user.Name))
	if nameLength < 3 {
		return errors.New("User name is too short.")
	}

	if nameLength > 100 {
		return errors.New("User name is too long.")
	}

	return nil
}

var ErrInvalidCredentials = errors.New("Invalid credentials.")

func (service *userService) verifyUser(name string, password string) (entity.User, error) {
	user, err := service.storage.findUserByName(name)
	if err != nil {
		service.hasher.verify(service.fake.Key, password)
		return user, ErrInvalidCredentials
	}

	isPasswordCorrect, newKeyChannel := service.hasher.verify(user.Key, password)
	if !isPasswordCorrect {
		return user, ErrInvalidCredentials
	}

	if newKeyChannel != nil {
		go service.updateUserKey(user, newKeyChannel)
	}

	return user, nil
}

func (service *userService) updateUserKey(user entity.User, newKeyChannel <-chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := service.storage.updateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}
