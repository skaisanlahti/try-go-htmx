package main

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/internal/auth"
	"github.com/skaisanlahti/try-go-htmx/internal/platform"
	"github.com/skaisanlahti/try-go-htmx/internal/todo"
	"github.com/skaisanlahti/try-go-htmx/internal/web"
)

func main() {
	settings := platform.ReadSettings("appsettings.json")
	database := platform.OpenDatabase(platform.DatabaseOptions{
		Driver:             "pgx",
		ConnectionString:   settings.Database.ConnectionString,
		MigrationDirectory: settings.Database.MigrationDirectory,
		MigrateOnStartup:   settings.Database.MigrateOnStartup,
	})

	server := platform.NewServer(settings.Address, database)
	sessionSecret := auth.NewSessionSecret(settings.Session.SecretLength)
	sessionDuration := time.Duration(settings.Session.SessionDurationMin * float64(time.Minute))
	sessionOptions := auth.SessionOptions{
		Secure:          settings.Session.Secure,
		CookieName:      settings.Session.CookieName,
		SessionSecret:   sessionSecret,
		SessionDuration: sessionDuration,
	}

	sessionStorage := auth.NewSessionStorage()
	sessionService := auth.NewSessionService(sessionOptions, sessionStorage)
	passwordOptions := auth.PasswordOptions{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: settings.Password.RecalculateOutdated,
	}

	passwordService := auth.NewPasswordService(passwordOptions)
	userStorage := auth.NewUserStorage(database)
	authService := auth.NewAuthenticationService(sessionService, passwordService, userStorage)
	authHtmxRenderer := auth.NewHtmxRenderer(web.TemplateFiles)
	authHtmxService := auth.NewHtmxService(authService, authHtmxRenderer)
	todoService := todo.NewTodoService(todo.NewTodoStorage(database))
	todoHtmxRenderer := todo.NewHtmxRenderer(web.TemplateFiles)
	todoHtmxService := todo.NewHtmxService(todoService, todoHtmxRenderer)
	middleware := platform.NewMiddlewareFactory(sessionService)
	web.MapAssets(server.Router)
	auth.MapRoutes(server.Router, authHtmxService, middleware)
	todo.MapRoutes(server.Router, todoHtmxService, middleware)
	server.Run()
}
