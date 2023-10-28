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

func (this *server) Run() {
	exit := make(chan struct{})
	go this.shutdown(exit)
	err := this.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-exit
}

func (this *server) shutdown(exit chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	token, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := this.Shutdown(token)
	if err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("Server closed.")
	err = this.database.Close()
	if err != nil {
		log.Fatalf("Database shutdown error: %v", err)
	}

	log.Println("Database closed.")
	log.Println("Shutdown finished.")
	close(exit)
}
