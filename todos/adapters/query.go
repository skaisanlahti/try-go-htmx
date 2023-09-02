package adapters

import (
	"database/sql"
	"log"
	"sync"
)

type Query struct {
	database *sql.DB
	queries  map[string]*sql.Stmt
	locker   sync.RWMutex
}

func NewQuery(database *sql.DB) *Query {
	return &Query{
		database: database,
		queries:  make(map[string]*sql.Stmt),
	}
}

func (this *Query) Prepare(query string) *sql.Stmt {
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
