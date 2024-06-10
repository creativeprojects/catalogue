package store

import (
	"errors"
	"sync"
)

// MemoryBucket represents a bucket in memory during a transaction
type MemoryBucket struct {
	data     map[string][]byte
	mutex    *sync.Mutex
	writable bool
}

// newMemoryBucket instantiate a new bucket in memory
func newMemoryBucket(data map[string][]byte, writable bool) *MemoryBucket {
	return &MemoryBucket{
		data:     data,
		mutex:    &sync.Mutex{},
		writable: writable,
	}
}

// Get returns a key from the bucket in memory
func (b *MemoryBucket) Get(key string) ([]byte, error) {
	if b == nil {
		return nil, ErrNullPointerBucket
	}
	if key == "" {
		return nil, ErrKeyNoName
	}

	var data []byte
	var ok bool

	b.mutex.Lock()
	defer b.mutex.Unlock()

	data, ok = b.data[key]
	if ok {
		return data, nil
	}
	return nil, ErrKeyNotFound
}

// Put sets a key in the memory bucket
func (b *MemoryBucket) Put(key string, data []byte) error {
	if b == nil {
		return ErrNullPointerBucket
	}
	if key == "" {
		return ErrKeyNoName
	}
	if !b.writable {
		return ErrBucketReadOnly
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.data[key] = data
	return nil
}

// Delete deletes a key from the memory bucket
func (b *MemoryBucket) Delete(key string) error {
	if b == nil {
		return ErrNullPointerBucket
	}
	if key == "" {
		return ErrKeyNoName
	}
	if !b.writable {
		return ErrBucketReadOnly
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, found := b.data[key]; found {
		delete(b.data, key)
	}
	return nil
}

func (b *MemoryBucket) CreateBucket(string) (Bucket, error) {
	return nil, errors.New("not implemented")
}
func (b *MemoryBucket) GetBucket(string) (Bucket, error) { return nil, errors.New("not implemented") }
func (b *MemoryBucket) DeleteBucket(string) error        { return errors.New("not implemented") }

// save returns the data in the bucket
func (b *MemoryBucket) save() map[string][]byte {
	return b.data
}

// Test the interface
var (
	_ Bucket = &MemoryBucket{}
)
