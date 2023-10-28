package entity

import "errors"

type User struct {
	Id   int
	Name string
	Key  []byte
}

func NewUser(name string, key []byte) User {
	return User{Name: name, Key: key}
}

func (this User) Validate() error {
	nameLength := len([]rune(this.Name))
	if nameLength < 3 {
		return errors.New("User name is too short.")
	}

	if nameLength > 100 {
		return errors.New("User name is too long.")
	}
	return nil
}
