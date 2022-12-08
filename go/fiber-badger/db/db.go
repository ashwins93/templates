package db

import "github.com/dgraph-io/badger/v3"

type Queries struct {
	db *badger.DB
}

func NewDb(db *badger.DB) *Queries {
	return &Queries{db}
}
