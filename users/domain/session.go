package domain

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/skaisanlahti/try-go-htmx/infrastructure"
)

const SessionDuration time.Duration = 60 * time.Minute
const SessionCookieName string = "sid"

type Session struct {
	Id      string
	UserId  int
	Expires time.Time
}

func NewSession(userId int) *Session {
	sessionId := uuid.New().String()
	expires := time.Now().Add(SessionDuration)
	return &Session{sessionId, userId, expires}
}

func NewSessionCookie(session *Session, mode string) *http.Cookie {
	maxAge := int(session.Expires.Sub(time.Now()).Seconds())
	secure := false
	if mode == infrastructure.ModeProduction {
		secure = true
	}

	return &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		Value:    session.Id,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   secure,
	}
}

func NewExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{Name: SessionCookieName,
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
