package domain

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Name     string
	Password []byte // hash
}

var (
	ErrUserNameTooLong    error = errors.New("User name is too long.")
	ErrPasswordTooLong    error = errors.New("Password is too long.")
	ErrPasswordHashFailed error = errors.New("Failed to hash password.")
)

func NewUser(name string, password string) (User, error) {
	if len([]rune(name)) > 100 {
		return User{}, ErrUserNameTooLong
	}

	if len(password) > 72 {
		return User{}, ErrPasswordTooLong
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return User{}, ErrPasswordHashFailed
	}

	return User{Name: name, Password: passwordHash}, nil
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}

	return hash, nil
}

func IsPasswordValid(hashedPassword []byte, candidatePassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, candidatePassword)
	return err == nil
}
