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
	sessionStorage := auth.NewSessionStorage()
	sessionService := auth.NewSessionService(auth.SessionOptions{
		Secure:          settings.Session.Secure,
		CookieName:      settings.Session.CookieName,
		SessionSecret:   sessionSecret,
		SessionDuration: sessionDuration,
		SessionStorage:  sessionStorage,
	})

	passwordService := auth.NewPasswordService(auth.PasswordOptions{
		Time:                settings.Password.Time,
		Memory:              settings.Password.Memory,
		Threads:             settings.Password.Threads,
		SaltLength:          settings.Password.SaltLength,
		KeyLength:           settings.Password.KeyLength,
		RecalculateOutdated: settings.Password.RecalculateOutdated,
	})

	authService := auth.NewAuthenticationService(sessionService, passwordService, auth.NewUserStorage(database))
	htmxAuthController := auth.NewHtmxController(authService, auth.NewHtmxRenderer(web.TemplateFiles))

	todoService := todo.NewTodoService(todo.NewTodoStorage(database))
	htmxTodoController := todo.NewHtmxController(todoService, todo.NewHtmxRenderer(web.TemplateFiles))

	middleware := platform.NewMiddlewareFactory(sessionService)

	web.MapAssets(server.Router)
	auth.MapRoutes(server.Router, htmxAuthController, middleware)
	todo.MapRoutes(server.Router, htmxTodoController, middleware)

	server.Run()
}
