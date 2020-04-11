package store

import (
	"sync"
)

type MemoryStore struct {
	buckets map[string]*MemoryBucket
	mutex   *sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets: make(map[string]*MemoryBucket, 0),
		mutex:   &sync.Mutex{},
	}
}

// GetBucket returns a bucket from its name. If it does not exists, a new empty bucket will be returned.
func (s *MemoryStore) GetBucket(bucket string) (Bucket, error) {
	if bucket == "" {
		return nil, ErrBucketNoName
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	b, ok := s.buckets[bucket]
	if !ok {
		b = NewMemoryBucket()
		s.buckets[bucket] = b
	}
	return b, nil
}

func (s *MemoryStore) GetKey(bucket string, key string) ([]byte, error) {
	b, err := s.GetBucket(bucket)
	if err != nil {
		return nil, err
	}
	return b.GetKey(key)
}

func (s *MemoryStore) SetKey(bucket string, key string, data []byte) error {
	b, err := s.GetBucket(bucket)
	if err != nil {
		return err
	}
	return b.SetKey(key, data)
}

func (s *MemoryStore) GetKeyUint64(bucket string, key string) (uint64, error) {
	b, err := s.GetBucket(bucket)
	if err != nil {
		return 0, err
	}
	return b.GetKeyUint64(key)
}

func (s *MemoryStore) SetKeyUint64(bucket string, key string, data uint64) error {
	b, err := s.GetBucket(bucket)
	if err != nil {
		return err
	}
	return b.SetKeyUint64(key, data)
}

// Begin a transaction
func (s *MemoryStore) Begin(writable bool) (Transaction, error) {
	return &MemoryTransaction{}, nil
}

// Test the interface
var (
	_ Store = &MemoryStore{}
)
