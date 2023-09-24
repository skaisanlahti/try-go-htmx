package settings

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func UseGracefulShutdown(server *http.Server, database *sql.DB) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Received an interrupt signal, shutting down...")
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		err = database.Close()
		if err != nil {
			log.Printf("Database shutdown error: %v", err)
		}

		os.Exit(0)
	}()

}
