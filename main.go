package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/assets"
	"github.com/skaisanlahti/try-go-htmx/settings"
	"github.com/skaisanlahti/try-go-htmx/todos"
	"github.com/skaisanlahti/try-go-htmx/users"
	"github.com/skaisanlahti/try-go-htmx/users/sessions"
)

func main() {
	variables := settings.ReadEnvironment(".env.development")
	database := settings.OpenDatabase(settings.DatabaseOptions{
		Driver:             "pgx",
		ConnectionString:   variables.Database,
		MigrationDirectory: "migrations/up",
		MigrateOnStartup:   settings.IsDevelopment(variables.Mode),
	})

	router := http.NewServeMux()
	session := sessions.NewStore(sessions.StoreOptions{
		CookieName:      "sid",
		SessionDuration: 60 * time.Minute,
		SessionSecret:   sessions.NewSecret(32),
		SessionStorage:  sessions.NewMemorySessionRepository(),
		Secure:          settings.IsProduction(variables.Mode),
	})

	assets.UseAssets(router)
	users.UseUserRoutes(router, database, session)
	todos.UseTodoRoutes(router, database, session)
	server := http.Server{
		Addr:         variables.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	settings.UseGracefulShutdown(&server, database)
	log.Printf("Starting a server at %s", variables.Address)
	log.Panic(server.ListenAndServe())
}
