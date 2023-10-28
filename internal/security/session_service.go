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
	options SessionOptions
	storage *sessionStorage
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

func (this *sessionService) verifySession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(this.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := this.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := this.storage.findSessionBySessionId(sessionId)
	if err != nil {
		return err
	}

	if session.Expires.Before(time.Now()) {
		this.storage.deleteSession(session)
		return errors.New("Session has expired.")
	}

	err = this.storage.updateSession(session.Extend(this.options.Duration))
	if err != nil {
		return err
	}

	signedSession, err := this.newSignature(session.Id)
	if err != nil {
		return err
	}

	newCookie := this.newSessionCookie(signedSession)
	http.SetCookie(response, newCookie)
	return nil
}

func (this *sessionService) sessionExists(request *http.Request) bool {
	cookie, err := request.Cookie(this.options.CookieName)
	if err != nil {
		return false
	}

	sessionId, err := this.verifySignature(cookie.Value)
	if err != nil {
		return false
	}

	_, err = this.storage.findSessionBySessionId(sessionId)
	if err != nil {
		return false
	}

	return true
}

func (this *sessionService) startSession(response http.ResponseWriter, userId int) error {
	session := entity.NewSession(userId, this.options.Duration)
	err := this.storage.insertSession(session)
	if err != nil {
		return err
	}

	signedSession, err := this.newSignature(session.Id)
	if err != nil {
		return err
	}

	cookie := this.newSessionCookie(signedSession)
	http.SetCookie(response, cookie)
	return nil
}

func (this *sessionService) clearSession(response http.ResponseWriter, request *http.Request) error {
	cookie, err := request.Cookie(this.options.CookieName)
	if err != nil {
		return err
	}

	sessionId, err := this.verifySignature(cookie.Value)
	if err != nil {
		return err
	}

	session, err := this.storage.findSessionBySessionId(sessionId)
	if err != nil {
		return err
	}

	err = this.storage.deleteSession(session)
	if err != nil {
		return err
	}

	http.SetCookie(response, this.clearSessionCookie())
	return nil
}

func (this *sessionService) newSessionCookie(signedSession string) *http.Cookie {
	return &http.Cookie{
		Name:     this.options.CookieName,
		Path:     "/",
		Value:    signedSession,
		MaxAge:   int(this.options.Duration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   this.options.Secure,
	}
}

func (this *sessionService) clearSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     this.options.CookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   this.options.Secure,
	}
}

func (this *sessionService) newSignature(sessionId string) (string, error) {
	code := hmac.New(sha256.New, []byte(this.options.Secret))
	code.Write([]byte(this.options.CookieName))
	code.Write([]byte(sessionId))
	signature := code.Sum(nil)
	signedSession := sessionId + "." + string(signature)
	encodedSession := base64.URLEncoding.EncodeToString([]byte(signedSession))
	if len(encodedSession) > 4096 {
		return "", errors.New("Cookie value too long.")
	}

	return encodedSession, nil
}

func (this *sessionService) verifySignature(encodedSession string) (string, error) {
	signedSession, err := base64.URLEncoding.DecodeString(encodedSession)
	if err != nil {
		return "", err
	}

	split := strings.SplitN(string(signedSession), ".", 2)
	sessionId := split[0]
	signature := split[1]
	code := hmac.New(sha256.New, []byte(this.options.Secret))
	code.Write([]byte(this.options.CookieName))
	code.Write([]byte(sessionId))
	expectedSignature := code.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", errors.New("Invalid signature.")
	}

	return sessionId, nil
}
