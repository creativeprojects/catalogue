package store

import (
	"encoding/binary"
	"sync"
)

// MemoryBucket represents a bucket in memory (it is *never* saved)
type MemoryBucket struct {
	readKeys    map[string][]byte
	writeKeys   map[string][]byte
	deletedKeys map[string]bool
	mutex       *sync.Mutex
	writable    bool
}

// newMemoryBucket instantiate a new bucket in memory
func newMemoryBucket(keys map[string][]byte, writable bool) *MemoryBucket {
	return &MemoryBucket{
		readKeys:    keys,
		writeKeys:   make(map[string][]byte, 0),
		deletedKeys: make(map[string]bool, 0),
		mutex:       &sync.Mutex{},
		writable:    writable,
	}
}

// GetKey returns a key from the bucket in memory
func (b *MemoryBucket) GetKey(key string) ([]byte, error) {
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

	// Was the key deleted?
	deleted, ok := b.deletedKeys[key]
	if ok && deleted {
		return nil, ErrKeyNotFound
	}

	// Was the key updated?
	data, ok = b.writeKeys[key]
	if ok {
		return data, nil
	}

	// Last resort: key was originally in the bucket
	data, ok = b.readKeys[key]
	if ok {
		return data, nil
	}
	return nil, ErrKeyNotFound
}

// SetKey sets a key in the memory bucket
func (b *MemoryBucket) SetKey(key string, data []byte) error {
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

	b.writeKeys[key] = data
	// In case the key was previously deleted
	if _, found := b.deletedKeys[key]; found {
		delete(b.deletedKeys, key)
	}
	return nil
}

// DeleteKey deletes a key from the memory bucket
func (b *MemoryBucket) DeleteKey(key string) error {
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

	b.deletedKeys[key] = true
	if _, found := b.writeKeys[key]; found {
		delete(b.writeKeys, key)
	}
	return nil
}

// GetKeyUint64 returns a uint64 value from the bucket in memory
func (b *MemoryBucket) GetKeyUint64(key string) (uint64, error) {
	value, err := b.GetKey(key)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(value), nil
}

// SetKeyUint64 sets a uint64 value in the memory bucket
func (b *MemoryBucket) SetKeyUint64(key string, data uint64) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, data)

	return b.SetKey(key, buffer)
}

// Test the interface
var (
	_ Bucket = &MemoryBucket{}
)
