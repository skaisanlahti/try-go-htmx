package auth

import (
	"database/sql"
	"time"

	"github.com/skaisanlahti/try-go-htmx/internal/platform"
)

type authModule struct {
	*htmxService
	*sessionService
}

func NewModule(
	password platform.PasswordSettings,
	session platform.SessionSettings,
	database *sql.DB,
) *authModule {
	secret := newSessionSecret(session.SecretLength)
	duration := time.Duration(session.SessionDurationMin * float64(time.Minute))
	sessionOptions := sessionOptions{
		secure:     session.Secure,
		cookieName: session.CookieName,
		secret:     secret,
		duration:   duration,
	}

	sessionStorage := newSessionStorage()
	sessionService := newSessionService(sessionOptions, sessionStorage)
	passwordOptions := passwordOptions{
		Time:                password.Time,
		Memory:              password.Memory,
		Threads:             password.Threads,
		SaltLength:          password.SaltLength,
		KeyLength:           password.KeyLength,
		RecalculateOutdated: password.RecalculateOutdated,
	}

	passwordHasher := newPasswordHasher(passwordOptions)
	userStorage := newUserStorage(database)
	userService := newUserService(userStorage, passwordHasher)
	authenticationService := newAuthenticationService(sessionService, userService)
	htmxRenderer := newHtmxRenderer(platform.TemplateFiles)
	htmxService := newHtmxService(authenticationService, htmxRenderer)
	return &authModule{htmxService, sessionService}
}
