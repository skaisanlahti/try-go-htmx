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
	settings := todoapp.ReadSettings("appsettings.json")
	database := psql.OpenDatabase(psql.DatabaseOptions{
		Driver:             "pgx",
		ConnectionString:   settings.Database.ConnectionString,
		MigrationDirectory: settings.Database.MigrationDirectory,
		MigrateOnStartup:   settings.Database.MigrateOnStartup,
	})

	userAccessor := psql.NewUserAccessor(database)
	todoAccessor := psql.NewTodoAccessor(database)
	sessionAccessor := mem.NewSessionAccessor()
	passwordEncoder := argon2.NewEncoder(argon2.EncoderOptions{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: settings.Password.RecalculateOutdated,
	})

	sessionManager := todoapp.NewSessionManager(todoapp.SessionOptions{
		Secure:          settings.Session.Secure,
		CookieName:      settings.Session.CookieName,
		SessionSecret:   todoapp.NewSessionSecret(settings.Session.SecretLength),
		SessionDuration: time.Duration(settings.Session.SessionDurationMin * float64(time.Minute)),
		SessionAccessor: sessionAccessor,
	})

	userAuthenticator := todoapp.NewUserAuthenticator(settings, sessionManager, passwordEncoder, userAccessor)
	todoWriter := todoapp.NewTodoWriter(settings, todoAccessor)
	htmxClient := htmx.NewClient(userAuthenticator, todoWriter)
	server := http.NewServer(settings.Address, database, htmxClient)
	server.Run()
}
