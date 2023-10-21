package platform

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	*http.Server
	Router   *http.ServeMux
	database *sql.DB
}

func NewServer(address string, database *sql.DB) *server {
	router := http.NewServeMux()
	listener := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &server{listener, router, database}
}

func (server *server) Run() {
	log.Printf("Server listening to %s", server.Addr)
	server.listenForInterrupt()
	log.Panic(server.ListenAndServe())
}

func (server *server) listenForInterrupt() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	go server.shutdown(interruptSignal)
}

func (server *server) shutdown(interruptSignal <-chan os.Signal) {
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
