package passwords

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

type Argon2idOptions struct {
	Time       uint32
	Memory     uint32
	Threads    uint8
	SaltLength uint32
	KeyLength  uint32
	Version    int
}

var DefaultArgon2idOptions Argon2idOptions = Argon2idOptions{
	Time:       14,
	Memory:     128 * 1024,
	Threads:    4,
	SaltLength: 32,
	KeyLength:  64,
	Version:    argon2.Version,
}

type Argon2Encoder struct {
	options Argon2idOptions
}

func NewArgon2Encoder(options Argon2idOptions) *Argon2Encoder {
	encoder := &Argon2Encoder{options}
	return encoder
}

func (encoder *Argon2Encoder) NewKey(password string) ([]byte, error) {
	salt := newSalt(encoder.options.SaltLength)

	reportProblems := monitorEncodingTime()
	key := argon2.IDKey([]byte(password), salt, encoder.options.Time, encoder.options.Memory, encoder.options.Threads, encoder.options.KeyLength)
	reportProblems()

	encodedKey := encoder.encodeKey(salt, key, encoder.options)
	return encodedKey, nil
}

func (encoder *Argon2Encoder) VerifyKey(encodedKey []byte, candidatePassword string, recalculateOutdatedKeys bool) (bool, chan []byte) {
	salt, key, options, err := encoder.decodeKey(encodedKey)
	if err != nil {
		return false, nil
	}

	reportProblems := monitorEncodingTime()
	candidateKey := argon2.IDKey([]byte(candidatePassword), salt, options.Time, options.Memory, options.Threads, options.KeyLength)
	isPasswordCorrect := subtle.ConstantTimeCompare(key, candidateKey) == 1
	reportProblems()

	optionsOutdated := encoder.areOptionsOutdated(options)
	var newKeyChannel chan []byte
	if isPasswordCorrect && optionsOutdated && recalculateOutdatedKeys {
		newKeyChannel = make(chan []byte)
		go encoder.recalculateKey(candidatePassword, newKeyChannel)
	}

	return isPasswordCorrect, newKeyChannel
}

func (encoder *Argon2Encoder) recalculateKey(password string, newKeyChannel chan<- []byte) {
	defer close(newKeyChannel)
	newKey, _ := encoder.NewKey(password)
	newKeyChannel <- newKey
}

func (encoder *Argon2Encoder) encodeKey(salt []byte, key []byte, options Argon2idOptions) []byte {
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

func (encoder *Argon2Encoder) decodeKey(encodedKey []byte) ([]byte, []byte, *Argon2idOptions, error) {
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

	options := &Argon2idOptions{}
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

func (encoder *Argon2Encoder) areOptionsOutdated(options *Argon2idOptions) bool {
	if encoder.options.Time != options.Time {
		return true
	}

	if encoder.options.Memory != options.Memory {
		return true
	}

	if encoder.options.Threads != options.Threads {
		return true
	}

	if encoder.options.SaltLength != options.SaltLength {
		return true
	}

	if encoder.options.KeyLength != options.KeyLength {
		return true
	}

	return false
}

func newSalt(size uint32) []byte {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		log.Panicln(err)
	}

	return salt
}
