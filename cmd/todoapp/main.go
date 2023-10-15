package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/todoapp"
	"github.com/skaisanlahti/try-go-htmx/todoapp/argon2"
	"github.com/skaisanlahti/try-go-htmx/todoapp/htmx"
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
		Time:        settings.Password.Time,
		Memory:      settings.Password.Memory,
		Threads:     settings.Password.Threads,
		SaltLength:  settings.Password.SaltLength,
		KeyLength:   settings.Password.KeyLength,
		Recalculate: settings.Password.RecalculateOutdated,
	})

	sessionManager := todoapp.NewSessionManager(todoapp.SessionOptions{
		CookieName:      settings.Session.CookieName,
		SessionDuration: time.Duration(settings.Session.SessionDurationMin * float64(time.Minute)),
		SessionSecret:   todoapp.NewSessionSecret(settings.Session.SecretLength),
		SessionAccessor: sessionAccessor,
		Secure:          isProduction(settings.Mode),
	})

	userAuthenticator := todoapp.NewUserAuthenticator(settings, sessionManager, passwordEncoder, userAccessor)
	todoWriter := todoapp.NewTodoWriter(settings, todoAccessor)
	htmxClient := htmx.NewClient(userAuthenticator, todoWriter)
	router := http.NewServeMux()
	htmx.AddRoutes(router, htmxClient)
	server := &http.Server{
		Addr:         settings.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	gracefulShutdown(server, database)
	log.Printf("Starting a server at %s", settings.Address)
	log.Panic(server.ListenAndServe())
}

const (
	ModeDevelopment string = "Development"
	ModeProduction  string = "Production"
)

func isDevelopment(mode string) bool {
	return mode == ModeDevelopment
}

func isProduction(mode string) bool {
	return mode == ModeProduction
}

func gracefulShutdown(server *http.Server, database *sql.DB) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Received an interrupt signal, shutting down...")
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		err = database.Close()
		if err != nil {
			log.Printf("Database shutdown error: %v", err)
		}

		os.Exit(0)
	}()

}
