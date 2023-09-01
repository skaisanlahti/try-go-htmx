package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/app"
	"github.com/skaisanlahti/try-go-htmx/assets"
	"github.com/skaisanlahti/try-go-htmx/tasks"
	"github.com/skaisanlahti/try-go-htmx/todos"
)

func main() {
	variables := app.ReadEnvironment(".env.development")
	database := app.OpenDatabase(variables.Database)
	defer database.Close()

	router := http.NewServeMux()
	assets.RegisterHandlers(router)
	todos.RegisterHandlers(router, database)
	tasks.RegisterHandlers(router, database)

	server := http.Server{
		Addr:         variables.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	log.Printf("Starting a server at %s", variables.Address)
	log.Fatal(server.ListenAndServe())
}
