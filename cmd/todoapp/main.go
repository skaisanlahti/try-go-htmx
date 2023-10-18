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

	server := http.NewServer(settings.Address, database)

	// services
	userStorage := psql.NewUserStorage(database)
	todoStorage := psql.NewTodoStorage(database)
	sessionStorage := mem.NewSessionStorage()
	passwordService := argon2.NewPasswordService(argon2.Options{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: settings.Password.RecalculateOutdated,
	})

	sessionService := todoapp.NewSessionService(todoapp.SessionOptions{
		Secure:          settings.Session.Secure,
		CookieName:      settings.Session.CookieName,
		SessionSecret:   todoapp.NewSessionSecret(settings.Session.SecretLength),
		SessionDuration: time.Duration(settings.Session.SessionDurationMin * float64(time.Minute)),
		SessionStorage:  sessionStorage,
	})

	authService := todoapp.NewAuthenticationService(sessionService, passwordService, userStorage)
	todoService := todoapp.NewTodoService(settings, todoStorage)

	// clients
	htmx.NewClient(authService, todoService, server.Router)

	server.Run()
}
