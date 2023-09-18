package passwords

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type BcryptEncoder struct {
	Cost int
}

func NewBcryptEncoder(cost int) *BcryptEncoder {
	return &BcryptEncoder{cost}
}

func (encoder *BcryptEncoder) NewKey(password string) ([]byte, error) {
	if len(password) > 72 {
		return nil, errors.New("Password is too long.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), encoder.Cost)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}

	return hash, nil
}

func (hasher *BcryptEncoder) VerifyKey(hashedPassword []byte, candidatePassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(candidatePassword))
	return err == nil, err
}
