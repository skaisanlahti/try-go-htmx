package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordOptions struct {
	Time                uint32
	Memory              uint32
	Threads             uint8
	SaltLength          uint32
	KeyLength           uint32
	RecalculateOutdated bool
	Version             uint32
}

type passwordService struct {
	options PasswordOptions
}

func NewPasswordService(options PasswordOptions) *passwordService {
	options.Version = argon2.Version
	service := &passwordService{options}
	return service
}

func newSalt(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}

func (service *passwordService) newKey(password string) ([]byte, error) {
	salt := newSalt(service.options.SaltLength)
	key := argon2.IDKey([]byte(password), salt, service.options.Time, service.options.Memory, service.options.Threads, service.options.KeyLength)
	encodedKey := service.encodeKey(salt, key, service.options)
	return encodedKey, nil
}

func (service *passwordService) verifyKey(encodedKey []byte, candidatePassword string) (bool, chan []byte) {
	salt, key, options, err := service.decodeKey(encodedKey)
	if err != nil {
		return false, nil
	}

	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	optionsOutdated := service.areOptionsOutdated(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && optionsOutdated && service.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go service.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (service *passwordService) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := service.newKey(password)
	newKeyChannel <- newKey
}

func (service *passwordService) encodeKey(salt []byte, key []byte, options PasswordOptions) []byte {
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedKey := base64.RawStdEncoding.EncodeToString(key)
	fullEncodedKey := []byte(fmt.Sprintf(
		"$argon2id$version=%d$time=%d,memory=%d,threads=%d$%s$%s",
		options.Version,
		options.Time,
		options.Memory,
		options.Threads,
		encodedSalt,
		encodedKey,
	))

	return fullEncodedKey
}

func (service *passwordService) decodeKey(encodedKey []byte) ([]byte, []byte, *PasswordOptions, error) {
	parts := strings.Split(string(encodedKey), "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("Invalid key.")
	}

	version := 0
	_, err := fmt.Sscanf(parts[2], "version=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, errors.New("Incompatible version.")
	}

	options := &PasswordOptions{}
	_, err = fmt.Sscanf(parts[3], "time=%d,memory=%d,threads=%d", &options.Time, &options.Memory, &options.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	options.SaltLength = uint32(len(salt))
	key, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	options.KeyLength = uint32(len(key))
	return salt, key, options, nil
}

func (service *passwordService) areOptionsOutdated(options *PasswordOptions) bool {
	if service.options.Time != options.Time {
		return true
	}

	if service.options.Memory != options.Memory {
		return true
	}

	if service.options.Threads != options.Threads {
		return true
	}

	if service.options.SaltLength != options.SaltLength {
		return true
	}

	if service.options.KeyLength != options.KeyLength {
		return true
	}

	return false
}
