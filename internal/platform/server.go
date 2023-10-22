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
	exit := make(chan struct{})
	go server.shutdown(exit)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-exit
}

func (server *server) shutdown(exit chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := server.Shutdown(context)
	if err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("Server closed.")
	err = server.database.Close()
	if err != nil {
		log.Fatalf("Database shutdown error: %v", err)
	}

	log.Println("Database closed.")
	log.Println("Shutdown finished.")
	close(exit)
}
