package store

import (
	"sync"
)

type MemoryStore struct {
	buckets      map[string]map[string][]byte
	bucketsMutex *sync.Mutex
	writeMutex   *sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets:      make(map[string]map[string][]byte, 0),
		bucketsMutex: &sync.Mutex{},
		writeMutex:   &sync.Mutex{},
	}
}

// Begin a transaction
func (s *MemoryStore) Begin(writable bool) (Transaction, error) {
	close := func() {}
	if writable {
		s.writeMutex.Lock()
		close = func() {
			s.writeMutex.Unlock()
		}
	}
	return newMemoryTransaction(s, writable, close), nil
}

// Close the store
func (s *MemoryStore) Close() {
	// Is there anything to do there?
}

// GetBucket returns a bucket from its name. If it does not exists, nil will be returned
func (s *MemoryStore) getBucket(bucket string) (map[string][]byte, error) {
	s.bucketsMutex.Lock()
	defer s.bucketsMutex.Unlock()

	b, ok := s.buckets[bucket]
	if ok {
		return b, nil
	}
	return nil, nil
}

// Test the interface
var (
	_ Store = &MemoryStore{}
)
