package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/assets"
	"github.com/skaisanlahti/try-go-htmx/infrastructure"
	"github.com/skaisanlahti/try-go-htmx/sessions"
	"github.com/skaisanlahti/try-go-htmx/todos"
	"github.com/skaisanlahti/try-go-htmx/users"
)

func main() {
	variables := infrastructure.ReadEnvironment(".env.development")
	database := infrastructure.OpenDatabase(variables.Database)
	defer database.Close()

	router := http.NewServeMux()
	session := sessions.NewStore(sessions.StoreOptions{
		CookieName:        "sid",
		SessionDuration:   60 * time.Second,
		SessionSecret:     sessions.NewSecret(),
		SessionRepository: sessions.NewInMemoryRepository(),
		Secure:            infrastructure.IsProduction(variables.Mode),
	})

	assets.MapAssetHandlers(router)
	users.MapHtmxHandlers(router, database, session)
	todos.MapHtmxHandlers(router, database, session)

	server := http.Server{
		Addr:         variables.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("Starting a server at %s", variables.Address)
	log.Panic(server.ListenAndServe())
}
