package domain

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Name     string
	Password []byte // hash
}

func NewUser(name string, password string) (User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return User{}, err
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
