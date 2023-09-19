package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

func NewSecret(size uint32) string {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(key)
}

type SessionStorage interface {
	Find(sessionId string) (*Session, error)
	Add(session *Session) error
	Update(session *Session) error
	Remove(sessionId string) error
}

type StoreOptions struct {
	CookieName      string
	SessionSecret   string
	SessionDuration time.Duration
	SessionStorage  SessionStorage
	Secure          bool
}

type Store struct {
	StoreOptions
}

func NewStore(options StoreOptions) *Store {
	store := &Store{options}
	return store
}

func (store *Store) Add(userId int) (*http.Cookie, error) {
	newSession := NewSession(userId, store.SessionDuration)
	err := store.SessionStorage.Add(newSession)
	if err != nil {
		return nil, err
	}

	cookie, err := store.newSessionCookie(newSession)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (store *Store) Remove(request *http.Request) (*http.Cookie, error) {
	cookie, err := request.Cookie(store.CookieName)
	if err != nil {
		return nil, ErrMissingSessionCookie
	}

	sessionId, err := store.decodeSession(cookie.Value)
	if err != nil {
		return nil, ErrInvalidSessionCookie
	}
	err = store.SessionStorage.Remove(sessionId)
	if err != nil {
		return nil, err
	}

	return store.newExpiredSessionCookie(), nil
}

var (
	ErrInvalidSessionCookie error = errors.New("Invalid session cookie.")
	ErrMissingSessionCookie error = errors.New("Session cookie missing.")
)

func (store *Store) Validate(request *http.Request) (*Session, error) {
	cookie, err := request.Cookie(store.CookieName)
	if err != nil {
		return nil, ErrMissingSessionCookie
	}

	encodedSession := cookie.Value
	if encodedSession == "" {
		return nil, ErrInvalidSessionCookie
	}

	sessionId, err := store.decodeSession(encodedSession)
	if err != nil {
		return nil, err
	}

	session, err := store.SessionStorage.Find(sessionId)
	if err != nil {
		return nil, ErrInvalidSessionCookie
	}

	if session.Expires.Before(time.Now()) {
		store.SessionStorage.Remove(session.Id)
		return nil, ErrInvalidSessionCookie
	}

	return session, nil
}

func (store *Store) Extend(session *Session) (*http.Cookie, error) {
	session.Expires = time.Now().Add(store.SessionDuration)
	err := store.SessionStorage.Update(session)
	if err != nil {
		return nil, err
	}

	cookie, err := store.newSessionCookie(session)
	if err != nil {
		return nil, err
	}

	return cookie, err
}

func (store *Store) newSessionCookie(session *Session) (*http.Cookie, error) {
	encodedSession, err := store.encodeSession(session.Id)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     store.CookieName,
		Path:     "/",
		Value:    encodedSession,
		MaxAge:   int(store.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   store.Secure,
	}, nil
}

func (store *Store) newExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     store.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   store.Secure,
	}
}

func (store *Store) encodeSession(sessionId string) (string, error) {
	code := hmac.New(sha256.New, []byte(store.SessionSecret))
	code.Write([]byte(store.CookieName))
	code.Write([]byte(sessionId))
	signature := code.Sum(nil)
	signedSession := sessionId + "." + string(signature)
	encodedSession := base64.URLEncoding.EncodeToString([]byte(signedSession))
	if len(encodedSession) > 4096 {
		return "", errors.New("Cookie value too long.")
	}

	return encodedSession, nil
}

func (store *Store) decodeSession(encodedSession string) (string, error) {
	signedSession, err := base64.URLEncoding.DecodeString(encodedSession)
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(signedSession), ".", 2)
	sessionId := split[0]
	signature := split[1]
	code := hmac.New(sha256.New, []byte(store.SessionSecret))
	code.Write([]byte(store.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return string(sessionId), nil
}