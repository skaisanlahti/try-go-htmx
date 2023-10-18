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

type SessionStorage interface {
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
	SessionStorage  SessionStorage
}

type SessionService struct {
	SessionOptions
}

func NewSessionService(o SessionOptions) *SessionService {
	return &SessionService{o}
}

func (service *SessionService) StartSession(response http.ResponseWriter, userId int) error {
	session := NewSession(userId, service.SessionDuration)
	err := service.SessionStorage.AddSession(session)
	if err != nil {
		return err
	}

	signedSession, err := service.newSignature(session.Id)
	if err != nil {
		return err
	}

	cookie, err := service.newSessionCookie(signedSession)
	if err != nil {
		return err
	}

	http.SetCookie(response, cookie)
	return nil
}

func (service *SessionService) StopSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	err = service.SessionStorage.RemoveSession(sessionId)
	if err != nil {
		return err
	}

	http.SetCookie(response, service.clearSessionCookie())
	return nil
}

func (service *SessionService) VerifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := service.SessionStorage.FindSession(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		service.SessionStorage.RemoveSession(session.Id)
		return errors.New("Session has expired.")
	}

	session.Expires = time.Now().Add(service.SessionDuration)
	err = service.SessionStorage.UpdateSession(session)
	if err != nil {
		return err
	}

	signedSession, err := service.newSignature(session.Id)
	if err != nil {
		return err
	}

	newCookie, err := service.newSessionCookie(signedSession)
	if err != nil {
		return err
	}

	http.SetCookie(response, newCookie)
	return nil
}

func (service *SessionService) newSessionCookie(signedSession string) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     service.CookieName,
		Path:     "/",
		Value:    signedSession,
		MaxAge:   int(service.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   service.Secure,
	}, nil
}

func (service *SessionService) clearSessionCookie() *http.Cookie {
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

func (service *SessionService) newSignature(sessionId string) (string, error) {
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

func (service *SessionService) verifySignature(encodedSession string) (string, error) {
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
