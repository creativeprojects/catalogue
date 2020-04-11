package database

import (
	"github.com/creativeprojects/catalogue/store"
)

type Database struct {
	storage store.Store
}

func NewDatabase(s store.Store) *Database {
	return &Database{
		storage: s,
	}
}
