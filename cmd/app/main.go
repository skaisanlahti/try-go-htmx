package main

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/internal/client/htmx"
	"github.com/skaisanlahti/try-go-htmx/internal/platform"
	"github.com/skaisanlahti/try-go-htmx/internal/security"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
)

func main() {
	// platform
	settings := platform.ReadSettings("appsettings.json")
	database := platform.OpenDatabase(platform.DatabaseOptions{
		Driver:             "pgx",
		ConnectionString:   settings.Database.ConnectionString,
		MigrationDirectory: settings.Database.MigrationDirectory,
		MigrateOnStartup:   settings.Database.MigrateOnStartup,
	})

	server := platform.NewServer(settings.Address, database)

	// services
	passwordOptions := security.PasswordOptions{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: true,
	}

	sessionOptions := security.SessionOptions{
		CookieName: settings.Session.CookieName,
		Secure:     settings.Session.Secure,
		Secret:     security.NewSessionSecret(settings.Session.SecretLength),
		Duration:   time.Duration(settings.Session.SessionDurationMin * float64(time.Minute)),
	}

	security := security.NewSecurityService(database, passwordOptions, sessionOptions)
	todo := todo.NewTodoService(database)

	// clients
	htmx.NewClient(security, todo, server.Router)

	// start
	server.Run()
}
