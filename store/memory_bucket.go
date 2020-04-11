package store

import (
	"encoding/binary"
	"errors"
	"sync"
)

// MemoryBucket represents a bucket in memory (it is *never* saved)
type MemoryBucket struct {
	keys  map[string][]byte
	mutex *sync.Mutex
}

// NewMemoryBucket instantiate a new bucket in memory
func NewMemoryBucket() *MemoryBucket {
	return &MemoryBucket{
		keys:  make(map[string][]byte, 0),
		mutex: &sync.Mutex{},
	}
}

// GetKey returns a key from the bucket in memory
func (b *MemoryBucket) GetKey(key string) ([]byte, error) {
	if b == nil {
		return nil, errors.New("Null pointer bucket")
	}
	if key == "" {
		return nil, ErrKeyNoName
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	data, ok := b.keys[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return data, nil
}

// SetKey sets a key in the memory bucket
func (b *MemoryBucket) SetKey(key string, data []byte) error {
	if b == nil {
		return errors.New("Null pointer bucket")
	}
	if key == "" {
		return ErrKeyNoName
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.keys[key] = data
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
