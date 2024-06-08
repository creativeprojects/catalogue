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
func (s *MemoryStore) Begin(writeable bool) (Transaction, error) {
	return s.begin(writeable)
}

func (s *MemoryStore) begin(writeable bool) (*MemoryTransaction, error) {
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

	if writeable {
		// Have a full lock on writeable buckets
		s.writeMutex.Lock()

		// Unlock at the end
		close = func() {
			s.writeMutex.Unlock()
			remove()
		}
	}

	tx := newMemoryTransaction(s, writeable, close)
	s.transactions[transactionID] = tx
	return tx, nil
}

// Close the store
func (s *MemoryStore) Close() {
	if len(s.transactions) > 0 {
		panic(fmt.Errorf("Store was closed with %d running transaction(s)", len(s.transactions)))
	}
}

func (s *MemoryStore) Update(job func(transaction Transaction) error) error {
	t, err := s.begin(true)
	if err != nil {
		return err
	}

	// Make sure the transaction rolls back in the event of a panic.
	defer t.Rollback()

	// If an error is returned from the function then rollback and return error.
	err = job(t)
	if err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}

func (s *MemoryStore) View(job func(transaction Transaction) error) error {
	t, err := s.begin(false)
	if err != nil {
		return err
	}

	// Make sure the transaction rolls back in the event of a panic.
	defer t.Rollback()

	// If an error is returned from the function then rollback and return error.
	err = job(t)
	if err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
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
