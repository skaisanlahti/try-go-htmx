package http

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skaisanlahti/try-go-htmx/todoapp/htmx"
)

type Server struct {
	*http.Server
	database *sql.DB
}

func NewServer(address string, database *sql.DB, htmxClient *htmx.Client) *Server {
	router := http.NewServeMux()
	htmx.AddRoutes(router, htmxClient)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{server, database}
}

func (server *Server) Run() {
	log.Printf("Server listening to %s", server.Addr)
	server.listenForInterrupt()
	log.Panic(server.ListenAndServe())
}

func (server *Server) listenForInterrupt() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	go server.shutdown(interruptSignal)
}

func (server *Server) shutdown(interruptSignal <-chan os.Signal) {
	<-interruptSignal
	log.Println("Received an interrupt signal, shutting down...")
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	err = server.database.Close()
	if err != nil {
		log.Printf("Database shutdown error: %v", err)
	}

	os.Exit(0)
}
