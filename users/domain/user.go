package domain

import (
	"errors"
)

type User struct {
	Id       int
	Name     string
	Password []byte
}

func NewUser(name string, key []byte) (User, error) {
	length := len([]rune(name))
	if length < 3 {
		return User{}, errors.New("User name is too short.")
	}

	if length > 100 {
		return User{}, errors.New("User name is too long.")
	}

	return User{Name: name, Password: key}, nil
}
