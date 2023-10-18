package main

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/todoapp"
	"github.com/skaisanlahti/try-go-htmx/todoapp/argon2"
	"github.com/skaisanlahti/try-go-htmx/todoapp/htmx"
	"github.com/skaisanlahti/try-go-htmx/todoapp/http"
	"github.com/skaisanlahti/try-go-htmx/todoapp/mem"
	"github.com/skaisanlahti/try-go-htmx/todoapp/psql"
)

func main() {
	// infrastructure
	settings := todoapp.ReadSettings("appsettings.json")
	database := psql.OpenDatabase(psql.DatabaseOptions{
		Driver:             "pgx",
		ConnectionString:   settings.Database.ConnectionString,
		MigrationDirectory: settings.Database.MigrationDirectory,
		MigrateOnStartup:   settings.Database.MigrateOnStartup,
	})

	server := http.CreateServer(settings.Address, database)

	// services
	userStorage := psql.CreateUserStorage(database)
	todoStorage := psql.CreateTodoStorage(database)
	sessionStorage := mem.CreateSessionStorage()
	passwordService := argon2.CreateService(argon2.Options{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: settings.Password.RecalculateOutdated,
	})

	sessionService := todoapp.CreateSessionService(todoapp.SessionOptions{
		Secure:          settings.Session.Secure,
		CookieName:      settings.Session.CookieName,
		SessionSecret:   todoapp.CreateSessionSecret(settings.Session.SecretLength),
		SessionDuration: time.Duration(settings.Session.SessionDurationMin * float64(time.Minute)),
		SessionStorage:  sessionStorage,
	})

	authService := todoapp.CreateAuthenticationService(sessionService, passwordService, userStorage)
	todoService := todoapp.CreateTodoService(settings, todoStorage)

	// clients
	htmx.CreateClient(authService, todoService, server.Router)

	server.Run()
}
