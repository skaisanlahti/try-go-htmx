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

func (service *SessionManager) StartSession(response http.ResponseWriter, userId int) error {
	session := NewSession(userId, service.SessionDuration)
	err := service.SessionAccessor.AddSession(session)
	if err != nil {
		return err
	}

	cookie, err := service.newSessionCookie(session)
	if err != nil {
		return err
	}

	http.SetCookie(response, cookie)
	return nil
}

func (service *SessionManager) StopSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.decodeSession(cookie.Value)
	if err != nil {
		return err
	}

	err = service.SessionAccessor.RemoveSession(sessionId)
	if err != nil {
		return err
	}

	http.SetCookie(response, service.newExpiredSessionCookie())
	return nil
}

func (service *SessionManager) VerifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.decodeSession(cookie.Value)
	if err != nil {
		return err
	}

	session, err := service.SessionAccessor.FindSession(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		service.SessionAccessor.RemoveSession(session.Id)
		return errors.New("Session has expired.")
	}

	session.Expires = time.Now().Add(service.SessionDuration)
	err = service.SessionAccessor.UpdateSession(session)
	if err != nil {
		return err
	}

	newCookie, err := service.newSessionCookie(session)
	if err != nil {
		return err
	}

	http.SetCookie(response, newCookie)
	return nil
}

func (service *SessionManager) newSessionCookie(session *Session) (*http.Cookie, error) {
	encodedSession, err := service.encodeSession(session.Id)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     service.CookieName,
		Path:     "/",
		Value:    encodedSession,
		MaxAge:   int(service.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   service.Secure,
	}, nil
}

func (service *SessionManager) newExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     service.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   service.Secure,
	}
}

func (service *SessionManager) encodeSession(sessionId string) (string, error) {
	code := hmac.New(sha256.New, []byte(service.SessionSecret))
	code.Write([]byte(service.CookieName))
	code.Write([]byte(sessionId))
	signature := code.Sum(nil)
	signedSession := sessionId + "." + string(signature)
	encodedSession := base64.URLEncoding.EncodeToString([]byte(signedSession))
	if len(encodedSession) > 4096 {
		return "", errors.New("Cookie value too long.")
	}

	return encodedSession, nil
}

func (service *SessionManager) decodeSession(encodedSession string) (string, error) {
	signedSession, err := base64.URLEncoding.DecodeString(encodedSession)
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(signedSession), ".", 2)
	sessionId := split[0]
	signature := split[1]
	code := hmac.New(sha256.New, []byte(service.SessionSecret))
	code.Write([]byte(service.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return sessionId, nil
}
