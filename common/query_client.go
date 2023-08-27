package common

import (
	"database/sql"
	"log"
	"sync"
)

type QueryClient struct {
	database *sql.DB
	queries  map[string]*sql.Stmt
	locker   sync.RWMutex
}

func NewQueryClient(database *sql.DB) *QueryClient {
	return &QueryClient{
		database: database,
		queries:  make(map[string]*sql.Stmt),
	}
}

func (this *QueryClient) Prepare(query string) *sql.Stmt {
	this.locker.RLock()
	prepared, ok := this.queries[query]
	this.locker.RUnlock()
	if ok {
		return prepared
	}

	prepared, err := this.database.Prepare(query)
	if err != nil {
		log.Printf("Failed to prepare query %s", query)
		log.Fatal(err)
	}

	this.locker.Lock()
	this.queries[query] = prepared
	this.locker.Unlock()
	return prepared
}
