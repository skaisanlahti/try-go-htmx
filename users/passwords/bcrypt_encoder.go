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

	reportProblems := monitorEncodingTime()
	key, err := bcrypt.GenerateFromPassword([]byte(password), encoder.Cost)
	reportProblems()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return key, nil
}

func (encoder *BcryptEncoder) VerifyKey(key []byte, candidatePassword string, recalculateOutdatedKeys bool) (bool, chan []byte) {
	reportProblems := monitorEncodingTime()
	isPasswordCorrect := bcrypt.CompareHashAndPassword(key, []byte(candidatePassword)) == nil
	reportProblems()

	costOutdated := encoder.isCostOutdated(key)
	var newKeyChannel chan []byte
	if isPasswordCorrect && costOutdated && recalculateOutdatedKeys {
		newKeyChannel = make(chan []byte)
		go encoder.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (encoder *BcryptEncoder) isCostOutdated(key []byte) bool {
	cost, err := bcrypt.Cost(key)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return encoder.Cost != cost
}

func (encoder *BcryptEncoder) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, err := encoder.NewKey(password)
	if err != nil {
		return
	}

	newKeyChannel <- newKey
}
