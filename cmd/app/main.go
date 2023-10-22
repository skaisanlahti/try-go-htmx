package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/try-go-htmx/internal/modules/auth"
	"github.com/skaisanlahti/try-go-htmx/internal/modules/todo"
	"github.com/skaisanlahti/try-go-htmx/internal/platform"
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

	authModule := auth.NewModule(settings.Password, settings.Session, database)
	todoModule := todo.NewModule(database)
	middleware := platform.NewMiddlewareFactory(authModule)

	platform.MapAssets(server.Router)
	auth.MapRoutes(server.Router, authModule, middleware)
	todo.MapRoutes(server.Router, todoModule, middleware)
	server.Run()
}
