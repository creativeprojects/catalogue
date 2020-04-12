package store

import (
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
)

// BoltBucket represents a bucket in memory (it is *never* saved)
type BoltBucket struct {
	bucket *bolt.Bucket
}

// newBoltBucket instantiate a new bucket in memory
func newBoltBucket(bucket *bolt.Bucket) *BoltBucket {
	return &BoltBucket{
		bucket: bucket,
	}
}

// GetKey returns a key from the bucket in memory
func (b *BoltBucket) GetKey(key string) ([]byte, error) {
	if b == nil {
		return nil, ErrNullPointerBucket
	}
	if key == "" {
		return nil, ErrKeyNoName
	}
	bucket := b.bucket.Get([]byte(key))
	if bucket == nil {
		return nil, ErrKeyNotFound
	}
	return bucket, nil
}

// SetKey sets a key in the memory bucket
func (b *BoltBucket) SetKey(key string, data []byte) error {
	if b == nil {
		return ErrNullPointerBucket
	}
	if key == "" {
		return ErrKeyNoName
	}
	if !b.bucket.Writable() {
		return ErrBucketReadOnly
	}

	return b.bucket.Put([]byte(key), data)
}

// DeleteKey deletes a key from the memory bucket
func (b *BoltBucket) DeleteKey(key string) error {
	if b == nil {
		return ErrNullPointerBucket
	}
	if key == "" {
		return ErrKeyNoName
	}
	if !b.bucket.Writable() {
		return ErrBucketReadOnly
	}

	return b.bucket.Delete([]byte(key))
}

// GetKeyUint64 returns a uint64 value from the bucket in memory
func (b *BoltBucket) GetKeyUint64(key string) (uint64, error) {
	value, err := b.GetKey(key)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(value), nil
}

// SetKeyUint64 sets a uint64 value in the memory bucket
func (b *BoltBucket) SetKeyUint64(key string, data uint64) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, data)

	return b.SetKey(key, buffer)
}

// Test the interface
var (
	_ Bucket = &BoltBucket{}
)
