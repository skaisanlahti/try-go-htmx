package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skaisanlahti/test-go/environment"
	"github.com/skaisanlahti/test-go/todos"
)

//go:embed public/*
var embeddedFiles embed.FS

const (
	embeddedFilesRoot = "public"
	assetsPath        = "/public/"
)

func mountAssetFiles() fs.FS {
	assetFiles, err := fs.Sub(embeddedFiles, embeddedFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	return assetFiles
}

func main() {
	variables := environment.Read(".env.development")
	database := environment.OpenDatabase(variables.Database)
	defer database.Close()

	assetFiles := mountAssetFiles()
	router := http.NewServeMux()
	router.Handle(assetsPath, http.StripPrefix(assetsPath, http.FileServer(http.FS(assetFiles))))
	todos.RegisterHandlers(router, database)

	log.Printf("Starting a server at %s", variables.Address)
	log.Fatal(http.ListenAndServe(variables.Address, router))
}
