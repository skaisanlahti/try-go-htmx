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

	fakeUser := entity.NewUser("username", fakeKey)
	return &userService{storage, password, fakeUser}
}

var ErrUserAlreadyExists = errors.New("Username already exists.")

func (this *userService) addUser(name string, password string) (int, error) {
	nameExists := this.storage.userExists(name)
	if nameExists {
		return 0, ErrUserAlreadyExists
	}

	key, err := this.hasher.hash(password)
	if err != nil {
		return 0, err
	}

	user := entity.NewUser(name, key)
	err = user.Validate()
	if err != nil {
		return 0, err
	}

	userId, err := this.storage.insertUser(user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

var ErrInvalidCredentials = errors.New("Invalid credentials.")

func (this *userService) verifyUser(name string, password string) (entity.User, error) {
	user, err := this.storage.findUserByName(name)
	if err != nil {
		this.hasher.verify(this.fake.Key, password)
		return user, ErrInvalidCredentials
	}

	isPasswordCorrect, newKeyChannel := this.hasher.verify(user.Key, password)
	if !isPasswordCorrect {
		return user, ErrInvalidCredentials
	}

	if newKeyChannel != nil {
		go this.updateUserKey(user, newKeyChannel)
	}

	return user, nil
}

func (this *userService) updateUserKey(user entity.User, newKeyChannel <-chan []byte) {
	newKey, ok := <-newKeyChannel
	if !ok {
		log.Printf("User: %s | Key update failed: recalculation failed.", user.Name)
		return
	}

	user.Key = newKey
	err := this.storage.updateUserKey(user)
	if err != nil {
		log.Printf("User: %s | Key update failed: database update failed.", user.Name)
	}
}
