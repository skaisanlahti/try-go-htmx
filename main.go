package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/assets"
	"github.com/skaisanlahti/try-go-htmx/infrastructure"
	"github.com/skaisanlahti/try-go-htmx/todos"
)

func main() {
	variables := infrastructure.ReadEnvironment(".env.development")
	database := infrastructure.OpenDatabase(variables.Database)
	defer database.Close()

	router := http.NewServeMux()
	assets.MapAssetHandlers(router)
	todos.MapHtmxHandlers(router, database)

	server := http.Server{
		Addr:         variables.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	log.Printf("Starting a server at %s", variables.Address)
	log.Panic(server.ListenAndServe())
}
