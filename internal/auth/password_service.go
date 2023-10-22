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

func NewPasswordService(options PasswordOptions) *passwordService {
	options.Version = argon2.Version
	service := &passwordService{options}
	return service
}

type passwordService struct {
	options PasswordOptions
}

func newSalt(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}

func (service *passwordService) hash(password string) ([]byte, error) {
	salt := newSalt(service.options.SaltLength)
	key := argon2.IDKey([]byte(password), salt, service.options.Time, service.options.Memory, service.options.Threads, service.options.KeyLength)
	encodedKey := service.encode(salt, key, service.options)
	return encodedKey, nil
}

func (service *passwordService) verify(encodedKey []byte, candidatePassword string) (bool, chan []byte) {
	salt, key, options, err := service.decode(encodedKey)
	if err != nil {
		return false, nil
	}

	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	ok := service.check(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && !ok && service.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go service.rehash(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (service *passwordService) rehash(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := service.hash(password)
	newKeyChannel <- newKey
}

func (service *passwordService) encode(salt []byte, key []byte, options PasswordOptions) []byte {
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

func (service *passwordService) decode(encodedKey []byte) ([]byte, []byte, *PasswordOptions, error) {
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

func (service *passwordService) check(options *PasswordOptions) bool {
	if service.options.Time != options.Time {
		return false
	}

	if service.options.Memory != options.Memory {
		return false
	}

	if service.options.Threads != options.Threads {
		return false
	}

	if service.options.SaltLength != options.SaltLength {
		return false
	}

	if service.options.KeyLength != options.KeyLength {
		return false
	}

	return true
}
