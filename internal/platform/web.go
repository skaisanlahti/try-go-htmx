package platform

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web/dist/*
var assetFiles embed.FS

const (
	assetFilesRoot = "web/dist"
	assetPath      = "/dist/"
)

func MapAssets(router *http.ServeMux) {
	assetFiles, err := fs.Sub(assetFiles, assetFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle(assetPath, http.StripPrefix(assetPath, http.FileServer(http.FS(assetFiles))))
}

//go:embed web/html/*.html
var TemplateFiles embed.FS
