package bcrypt

import (
	"errors"
	"log"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	options todoapp.PasswordSettings
}

func CreateService(options todoapp.PasswordSettings) *Service {
	return &Service{options}
}

func (service *Service) NewKey(password string) ([]byte, error) {
	if len(password) > 72 {
		return nil, errors.New("Password is too long.")
	}

	reportProblems := todoapp.MonitorEncodingTime()
	key, err := bcrypt.GenerateFromPassword([]byte(password), service.options.Cost)
	reportProblems()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return key, nil
}

func (service *Service) VerifyKey(key []byte, candidatePassword string) (bool, chan []byte) {
	reportProblems := todoapp.MonitorEncodingTime()
	isPasswordCorrect := bcrypt.CompareHashAndPassword(key, []byte(candidatePassword)) == nil
	reportProblems()

	costOutdated := service.isCostOutdated(key)
	var newKeyChannel chan []byte
	if isPasswordCorrect && costOutdated && service.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go service.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (service *Service) isCostOutdated(key []byte) bool {
	cost, err := bcrypt.Cost(key)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return service.options.Cost != cost
}

func (service *Service) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, err := service.NewKey(password)
	if err != nil {
		return
	}

	newKeyChannel <- newKey
}
