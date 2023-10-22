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

type passwordOptions struct {
	Time                uint32
	Memory              uint32
	Threads             uint8
	SaltLength          uint32
	KeyLength           uint32
	RecalculateOutdated bool
	Version             uint32
}

type passwordHasher struct {
	options passwordOptions
}

func newPasswordHasher(options passwordOptions) *passwordHasher {
	options.Version = argon2.Version
	return &passwordHasher{options}
}

func newSalt(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}

func (hasher *passwordHasher) hash(password string) ([]byte, error) {
	salt := newSalt(hasher.options.SaltLength)
	key := argon2.IDKey([]byte(password), salt, hasher.options.Time, hasher.options.Memory, hasher.options.Threads, hasher.options.KeyLength)
	encodedKey := hasher.encode(salt, key, hasher.options)
	return encodedKey, nil
}

func (hasher *passwordHasher) verify(encodedKey []byte, candidatePassword string) (bool, chan []byte) {
	salt, key, options, err := hasher.decode(encodedKey)
	if err != nil {
		return false, nil
	}

	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	ok := hasher.check(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && !ok && hasher.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go hasher.rehash(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (hasher *passwordHasher) rehash(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := hasher.hash(password)
	newKeyChannel <- newKey
}

func (hasher *passwordHasher) encode(salt []byte, key []byte, options passwordOptions) []byte {
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

func (hasher *passwordHasher) decode(encodedKey []byte) ([]byte, []byte, *passwordOptions, error) {
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
		return nil, nil, nil, errors.New("Incompatible argon2 version.")
	}

	options := &passwordOptions{}
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

func (hasher *passwordHasher) check(options *passwordOptions) bool {
	if hasher.options.Time != options.Time {
		return false
	}

	if hasher.options.Memory != options.Memory {
		return false
	}

	if hasher.options.Threads != options.Threads {
		return false
	}

	if hasher.options.SaltLength != options.SaltLength {
		return false
	}

	if hasher.options.KeyLength != options.KeyLength {
		return false
	}

	return true
}
