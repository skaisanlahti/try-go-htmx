package assets

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var embeddedFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

const (
	embeddedFilesRoot = "dist"
	path              = "/dist/"
)

func RegisterHandlers(router *http.ServeMux) {
	assetFiles, err := fs.Sub(embeddedFiles, embeddedFilesRoot)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle(path, http.StripPrefix(path, http.FileServer(http.FS(assetFiles))))
}

type HtmlTemplates struct {
	TodoPageTemplate *template.Template
}

func GetTemplates() *HtmlTemplates {
	todoPage := template.Must(template.ParseFS(templateFiles, "templates/todo_page.html"))
	return &HtmlTemplates{todoPage}
}
