package auth

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

type SessionOptions struct {
	Secure          bool
	CookieName      string
	SessionSecret   string
	SessionDuration time.Duration
}

func NewSessionService(options SessionOptions, storage *sessionStorage) *sessionService {
	return &sessionService{options, storage}
}

func NewSessionSecret(length uint32) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

type sessionService struct {
	options        SessionOptions
	sessionStorage *sessionStorage
}

func (service *sessionService) VerifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := service.sessionStorage.findSession(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		service.sessionStorage.removeSession(session.Id)
		return errors.New("Session has expired.")
	}

	session.Expires = time.Now().Add(service.options.SessionDuration)
	err = service.sessionStorage.updateSession(session)
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

func (service *sessionService) startSession(response http.ResponseWriter, userId int) error {
	session := newSession(userId, service.options.SessionDuration)
	err := service.sessionStorage.addSession(session)
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

func (service *sessionService) stopSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	err = service.sessionStorage.removeSession(sessionId)
	if err != nil {
		return err
	}

	http.SetCookie(response, service.clearSessionCookie())
	return nil
}

func (service *sessionService) newSessionCookie(signedSession string) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     service.options.CookieName,
		Path:     "/",
		Value:    signedSession,
		MaxAge:   int(service.options.SessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   service.options.Secure,
	}, nil
}

func (service *sessionService) clearSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     service.options.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   service.options.Secure,
	}
}

func (service *sessionService) newSignature(sessionId string) (string, error) {
	code := hmac.New(sha256.New, []byte(service.options.SessionSecret))
	code.Write([]byte(service.options.CookieName))
	code.Write([]byte(sessionId))
	signature := code.Sum(nil)
	signedSession := sessionId + "." + string(signature)
	encodedSession := base64.URLEncoding.EncodeToString([]byte(signedSession))
	if len(encodedSession) > 4096 {
		return "", errors.New("Cookie value too long.")
	}

	return encodedSession, nil
}

func (service *sessionService) verifySignature(encodedSession string) (string, error) {
	signedSession, err := base64.URLEncoding.DecodeString(encodedSession)
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(signedSession), ".", 2)
	sessionId := split[0]
	signature := split[1]
	code := hmac.New(sha256.New, []byte(service.options.SessionSecret))
	code.Write([]byte(service.options.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return sessionId, nil
}
