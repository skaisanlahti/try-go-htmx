package bcrypt

import (
	"errors"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
	"golang.org/x/crypto/bcrypt"
)

type Encoder struct {
	options todoapp.PasswordSettings
}

func NewEncoder(options todoapp.PasswordSettings) *Encoder {
	return &Encoder{options}
}

func (encoder *Encoder) NewKey(password string) ([]byte, error) {
	if len(password) > 72 {
		return nil, errors.New("Password is too long.")
	}

	reportProblems := todoapp.MonitorEncodingTime()
	key, err := bcrypt.GenerateFromPassword([]byte(password), encoder.options.Cost)
	reportProblems()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return key, nil
}

func (encoder *Encoder) VerifyKey(key []byte, candidatePassword string) (bool, chan []byte) {
	reportProblems := todoapp.MonitorEncodingTime()
	isPasswordCorrect := bcrypt.CompareHashAndPassword(key, []byte(candidatePassword)) == nil
	reportProblems()

	costOutdated := encoder.isCostOutdated(key)
	var newKeyChannel chan []byte
	if isPasswordCorrect && costOutdated && encoder.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go encoder.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (encoder *Encoder) isCostOutdated(key []byte) bool {
	cost, err := bcrypt.Cost(key)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return encoder.options.Cost != cost
}

func (encoder *Encoder) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, err := encoder.NewKey(password)
	if err != nil {
		return
	}

	newKeyChannel <- newKey
}
