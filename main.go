package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/skaisanlahti/test-go/environment"
	"github.com/skaisanlahti/test-go/todos"
)

func duration(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before := time.Now().UnixMilli()
		next(w, r)
		after := time.Now().UnixMilli()
		duration := after - before
		log.Printf("Request took %d ms", duration)
	})
}

//go:embed public/*
var embeddedFS embed.FS

const (
	embeddedFSRoot = "public"
	assetsPath     = "/public/"
)

func main() {
	environment := environment.New(".env.development")
	assetsFS, _ := fs.Sub(embeddedFS, embeddedFSRoot)
	http.Handle(assetsPath, http.StripPrefix(assetsPath, http.FileServer(http.FS(assetsFS))))

	service := todos.NewService()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, todos.HomePath, http.StatusTemporaryRedirect)
	})
	http.HandleFunc(todos.HomePath, duration(service.Home))
	http.HandleFunc(todos.GetTodosPath, duration(service.GetTodos))
	http.HandleFunc(todos.AddTodoPath, duration(service.AddTodo))
	http.HandleFunc(todos.RemoveTodoPath, duration(service.RemoveTodo))
	http.HandleFunc(todos.ToggleTodoPath, duration(service.ToggleTodo))

	log.Printf("Starting a server at %s", environment.Address)
	log.Fatal(http.ListenAndServe(environment.Address, nil))
}
