package todoapp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id      string
	UserId  int
	Expires time.Time
}

func NewSession(userId int, duration time.Duration) *Session {
	sessionId := uuid.New().String()
	expires := time.Now().Add(duration)
	return &Session{sessionId, userId, expires}
}

func NewSessionSecret(length uint32) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

type SessionAccessor interface {
	FindSession(sessionId string) (*Session, error)
	AddSession(session *Session) error
	UpdateSession(session *Session) error
	RemoveSession(sessionId string) error
}

type SessionOptions struct {
	Secure          bool
	CookieName      string
	SessionSecret   string
	SessionDuration time.Duration
	SessionAccessor SessionAccessor
}

type SessionManager struct {
	SessionOptions
}

func NewSessionManager(o SessionOptions) *SessionManager {
	return &SessionManager{o}
}

func (manager *SessionManager) StartSession(response http.ResponseWriter, userId int) error {
	session := NewSession(userId, manager.SessionDuration)
	err := manager.SessionAccessor.AddSession(session)
	if err != nil {
		return err
	}

	cookie, err := manager.newSessionCookie(session)
	if err != nil {
		return err
	}

	http.SetCookie(response, cookie)
	return nil
}

func (manager *SessionManager) StopSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(manager.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := manager.decodeSession(cookie.Value)
	if err != nil {
		return err
	}

	err = manager.SessionAccessor.RemoveSession(sessionId)
	if err != nil {
		return err
	}

	http.SetCookie(response, manager.newExpiredSessionCookie())
	return nil
}

func (manager *SessionManager) VerifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(manager.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := manager.decodeSession(cookie.Value)
	if err != nil {
		return err
	}

	session, err := manager.SessionAccessor.FindSession(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		manager.SessionAccessor.RemoveSession(session.Id)
		return errors.New("Session has expired.")
	}

	session.Expires = time.Now().Add(manager.SessionDuration)
	err = manager.SessionAccessor.UpdateSession(session)
	if err != nil {
		return err
	}

	newCookie, err := manager.newSessionCookie(session)
	if err != nil {
		return err
	}

	http.SetCookie(response, newCookie)
	return nil
}

func (manager *SessionManager) newSessionCookie(session *Session) (*http.Cookie, error) {
	encodedSession, err := manager.encodeSession(session.Id)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     manager.CookieName,
		Path:     "/",
		Value:    encodedSession,
		MaxAge:   int(manager.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   manager.Secure,
	}, nil
}

func (manager *SessionManager) newExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     manager.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   manager.Secure,
	}
}

func (manager *SessionManager) encodeSession(sessionId string) (string, error) {
	code := hmac.New(sha256.New, []byte(manager.SessionSecret))
	code.Write([]byte(manager.CookieName))
	code.Write([]byte(sessionId))
	signature := code.Sum(nil)
	signedSession := sessionId + "." + string(signature)
	encodedSession := base64.URLEncoding.EncodeToString([]byte(signedSession))
	if len(encodedSession) > 4096 {
		return "", errors.New("Cookie value too long.")
	}

	return encodedSession, nil
}

func (manager *SessionManager) decodeSession(encodedSession string) (string, error) {
	signedSession, err := base64.URLEncoding.DecodeString(encodedSession)
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(signedSession), ".", 2)
	sessionId := split[0]
	signature := split[1]
	code := hmac.New(sha256.New, []byte(manager.SessionSecret))
	code.Write([]byte(manager.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return sessionId, nil
}
