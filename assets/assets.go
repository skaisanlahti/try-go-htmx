package assets

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var embeddedFiles embed.FS

const (
	embeddedFilesRoot = "dist"
	path              = "/dist/"
)

func MapAssetHandlers(router *http.ServeMux) {
	assetFiles, err := fs.Sub(embeddedFiles, embeddedFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle(path, http.StripPrefix(path, http.FileServer(http.FS(assetFiles))))
}
