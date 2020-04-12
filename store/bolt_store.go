package store

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	db *bolt.DB
}

func NewBoltStore(database string) (*BoltStore, error) {
	db, err := bolt.Open(database, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot open database file '%s'", database)
	}
	return &BoltStore{
		db: db,
	}, nil
}

// Begin a transaction
func (s *BoltStore) Begin(writable bool) (Transaction, error) {
	tx, err := s.db.Begin(writable)
	if err != nil {
		return nil, err
	}
	return newBoltTransaction(tx), nil
}

// Close the store
func (s *BoltStore) Close() {
	s.db.Close()
}

// Test the interface
var (
	_ Store = &BoltStore{}
)
