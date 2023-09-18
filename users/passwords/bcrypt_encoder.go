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

	key, err := bcrypt.GenerateFromPassword([]byte(password), encoder.Cost)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}

	return key, nil
}

func (encoder *BcryptEncoder) VerifyKey(key []byte, candidatePassword string, recalculateOutdatedKeys bool) (bool, []byte, error) {
	err := bcrypt.CompareHashAndPassword(key, []byte(candidatePassword))
	isPasswordCorrect := err == nil
	return isPasswordCorrect, nil, err
}

func (encoder *BcryptEncoder) VerifyOptions(cost int) bool {
	return encoder.Cost == cost
}
