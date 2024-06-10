package store

import (
	"sync"
)

// MemoryBucket represents a bucket in memory during a transaction
type MemoryBucket struct {
	name  string
	data  map[string][]byte
	mutex *sync.Mutex
	tx    *MemoryTransaction
}

// newMemoryBucket instantiate a new bucket in memory
func newMemoryBucket(name string, data map[string][]byte, tx *MemoryTransaction) *MemoryBucket {
	return &MemoryBucket{
		name:  name,
		data:  data,
		mutex: &sync.Mutex{},
		tx:    tx,
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
	if !b.tx.IsWritable() {
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
	if !b.tx.IsWritable() {
		return ErrBucketReadOnly
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, found := b.data[key]; found {
		delete(b.data, key)
	}
	return nil
}

func (b *MemoryBucket) CreateBucket(name string) (Bucket, error) {
	return b.tx.CreateBucket(b.bucketName(name))
}

func (b *MemoryBucket) GetBucket(name string) (Bucket, error) {
	return b.tx.GetBucket(b.bucketName(name))
}

func (b *MemoryBucket) DeleteBucket(name string) error {
	return b.tx.DeleteBucket(b.bucketName(name))
}

func (b *MemoryBucket) bucketName(name string) string {
	return b.name + "/" + name
}

// getData returns the data in the bucket
func (b *MemoryBucket) getData() map[string][]byte {
	return b.data
}

// Test the interface
var (
	_ Bucket = &MemoryBucket{}
)
