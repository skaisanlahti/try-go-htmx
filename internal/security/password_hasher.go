package security

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

type PasswordHasher struct {
	options PasswordOptions
}

func NewPasswordHasher(options PasswordOptions) *PasswordHasher {
	options.Version = argon2.Version
	return &PasswordHasher{options}
}

func newSalt(length uint32) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Panicln(err)
	}

	return bytes
}

func (this *PasswordHasher) Hash(password string) ([]byte, error) {
	salt := newSalt(this.options.SaltLength)
	key := argon2.IDKey([]byte(password), salt, this.options.Time, this.options.Memory, this.options.Threads, this.options.KeyLength)
	encodedKey := this.encode(salt, key, this.options)
	return encodedKey, nil
}

func (this *PasswordHasher) Verify(encodedKey []byte, candidatePassword string) (bool, chan []byte) {
	salt, key, options, err := this.decode(encodedKey)
	if err != nil {
		return false, nil
	}

	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	ok := this.check(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && !ok && this.options.RecalculateOutdated {
		newKeyChannel = make(chan []byte)
		go this.rehash(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (this *PasswordHasher) rehash(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := this.Hash(password)
	newKeyChannel <- newKey
}

func (this *PasswordHasher) encode(salt []byte, key []byte, options PasswordOptions) []byte {
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

func (this *PasswordHasher) decode(encodedKey []byte) ([]byte, []byte, *PasswordOptions, error) {
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

func (this *PasswordHasher) check(options *PasswordOptions) bool {
	if this.options.Time != options.Time {
		return false
	}

	if this.options.Memory != options.Memory {
		return false
	}

	if this.options.Threads != options.Threads {
		return false
	}

	if this.options.SaltLength != options.SaltLength {
		return false
	}

	if this.options.KeyLength != options.KeyLength {
		return false
	}

	return true
}
