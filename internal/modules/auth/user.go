package auth

import (
	"errors"
)

type user struct {
	Id   int
	Name string
	Key  []byte
}

func newUser(name string, key []byte) (user, error) {
	nameLength := len([]rune(name))
	if nameLength < 3 {
		return user{}, errors.New("User name is too short.")
	}

	if nameLength > 100 {
		return user{}, errors.New("User name is too long.")
	}

	return user{Name: name, Key: key}, nil
}
