package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/assets"
	"github.com/skaisanlahti/try-go-htmx/infrastructure"
	"github.com/skaisanlahti/try-go-htmx/todos"
	"github.com/skaisanlahti/try-go-htmx/users"
	"github.com/skaisanlahti/try-go-htmx/users/memory"
)

func main() {
	variables := infrastructure.ReadEnvironment(".env.development")
	database := infrastructure.OpenDatabase(variables.Database)
	defer database.Close()

	sessions := memory.NewSessionStore()
	router := http.NewServeMux()
	assets.MapAssetHandlers(router)
	users.MapHtmxHandlers(router, database, sessions, variables.Mode)
	todos.MapHtmxHandlers(router, database, sessions, variables.Mode)

	server := http.Server{
		Addr:         variables.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("Starting a server at %s", variables.Address)
	log.Panic(server.ListenAndServe())
}
