package store

import (
	"fmt"
	"sync"
)

type MemoryStore struct {
	buckets           map[string]map[string][]byte
	bucketsMutex      *sync.Mutex
	writeMutex        *sync.Mutex
	transactions      map[uint]*MemoryTransaction
	transactionsMutex *sync.Mutex
	nextTransactionID uint
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets:           make(map[string]map[string][]byte, 0),
		bucketsMutex:      &sync.Mutex{},
		writeMutex:        &sync.Mutex{},
		transactions:      make(map[uint]*MemoryTransaction, 0),
		transactionsMutex: &sync.Mutex{},
	}
}

// Begin a transaction
func (s *MemoryStore) Begin(writable bool) (Transaction, error) {
	s.transactionsMutex.Lock()
	defer s.transactionsMutex.Unlock()

	s.nextTransactionID++
	transactionID := s.nextTransactionID

	// Callback to remove the transaction from the map once it's finished
	remove := func() {
		s.transactionsMutex.Lock()
		defer s.transactionsMutex.Unlock()

		delete(s.transactions, transactionID)
	}

	close := func() {
		remove()
	}

	if writable {
		// Have a full lock on writable buckets
		s.writeMutex.Lock()

		// Unlock at the end
		close = func() {
			s.writeMutex.Unlock()
			remove()
		}
	}

	tx := newMemoryTransaction(s, writable, close)
	s.transactions[transactionID] = tx
	return tx, nil
}

// Close the store
func (s *MemoryStore) Close() {
	if len(s.transactions) > 0 {
		panic(fmt.Errorf("Store was closed with %d running transaction(s)", len(s.transactions)))
	}
}

// getBucket returns a bucket from its name. If it does not exists, a new bucket will be created
func (s *MemoryStore) getBucket(bucket string) map[string][]byte {
	s.bucketsMutex.Lock()
	defer s.bucketsMutex.Unlock()

	b, ok := s.buckets[bucket]
	if !ok {
		b = make(map[string][]byte, 0)
		s.buckets[bucket] = b
	}
	return b
}

// deleteBucket removes the bucket from memory. It does not fail if the bucket was not found
func (s *MemoryStore) deleteBucket(bucket string) {
	s.bucketsMutex.Lock()
	defer s.bucketsMutex.Unlock()

	_, found := s.buckets[bucket]
	if found {
		delete(s.buckets, bucket)
	}
}

// Test the interface
var (
	_ Store = &MemoryStore{}
)
