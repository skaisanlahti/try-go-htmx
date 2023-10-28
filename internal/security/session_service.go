package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/entity"
)

type SessionOptions struct {
	Secure     bool
	CookieName string
	Secret     string
	Duration   time.Duration
}

type sessionService struct {
	options        SessionOptions
	sessionStorage *sessionStorage
}

func newSessionService(options SessionOptions, storage *sessionStorage) *sessionService {
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

func (service *sessionService) verifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := service.sessionStorage.findSessionBySessionId(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		service.sessionStorage.deleteSession(session)
		return errors.New("Session has expired.")
	}

	session.Expires = time.Now().Add(service.options.Duration)
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

func (service *sessionService) sessionExists(request *http.Request) bool {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return false
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return false
	}

	_, err = service.sessionStorage.findSessionBySessionId(sessionId)
	if err != nil {
		return false
	}

	return true
}

func (service *sessionService) startSession(response http.ResponseWriter, userId int) error {
	session := entity.NewSession(userId, service.options.Duration)
	err := service.sessionStorage.insertSession(session)
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

func (service *sessionService) clearSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := service.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := service.sessionStorage.findSessionBySessionId(sessionId)
	if err != nil {
		return err
	}

	err = service.sessionStorage.deleteSession(session)
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
		MaxAge:   int(service.options.Duration.Seconds()),
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
	code := hmac.New(sha256.New, []byte(service.options.Secret))
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
	code := hmac.New(sha256.New, []byte(service.options.Secret))
	code.Write([]byte(service.options.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return sessionId, nil
}
