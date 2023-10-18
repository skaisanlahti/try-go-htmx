package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/skaisanlahti/try-go-htmx/todoapp"
	"golang.org/x/crypto/argon2"
)

type Options struct {
	Time                uint32
	Memory              uint32
	Threads             uint8
	SaltLength          uint32
	KeyLength           uint32
	RecalculateOutdated bool
	Version             uint32
}

type PasswordService struct {
	options Options
}

func NewPasswordService(options Options) *PasswordService {
	options.Version = argon2.Version
	service := &PasswordService{options}
	return service
}

func (service *PasswordService) NewKey(password string) ([]byte, error) {
	salt := newSalt(service.options.SaltLength)

	reportProblems := todoapp.MonitorEncodingTime()
	key := argon2.IDKey([]byte(password), salt, service.options.Time, service.options.Memory, service.options.Threads, service.options.KeyLength)
	reportProblems()

	encodedKey := service.encodeKey(salt, key, service.options)
	return encodedKey, nil
}

func (service *PasswordService) VerifyKey(encodedKey []byte, candidatePassword string) (bool, chan []byte) {
	salt, key, options, err := service.decodeKey(encodedKey)
	if err != nil {
		return false, nil
	}

	reportProblems := todoapp.MonitorEncodingTime()
	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	reportProblems()

	optionsOutdated := service.areOptionsOutdated(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && optionsOutdated && service.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go service.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (service *PasswordService) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := service.NewKey(password)
	newKeyChannel <- newKey
}

func (service *PasswordService) encodeKey(salt []byte, key []byte, options Options) []byte {
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

func (service *PasswordService) decodeKey(encodedKey []byte) ([]byte, []byte, *Options, error) {
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

	options := &Options{}
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

func (service *PasswordService) areOptionsOutdated(options *Options) bool {
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

func newSalt(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}
