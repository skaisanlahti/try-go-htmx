package htmx

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web/dist/*
var assetFiles embed.FS

const (
	assetFilesRoot = "web/dist"    // Root of asset files in the embedded file system.
	assetPath      = "/htmx/dist/" // Asset path in request URL. Used in HTML templates as prefix for the asset files.
)

func (this *controller) getAssets(response http.ResponseWriter, request *http.Request) {
	// strip the file system path from assets
	assets, err := fs.Sub(assetFiles, assetFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	// strip request url from assets so that requests are properly mapped to the file system
	http.StripPrefix(assetPath, http.FileServer(http.FS(assets))).ServeHTTP(response, request)
}
