package store

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

type BoltBucket struct {
	bucket *bolt.Bucket
}

func newBoltBucket(bucket *bolt.Bucket) *BoltBucket {
	return &BoltBucket{
		bucket: bucket,
	}
}

func (b *BoltBucket) Get(key string) ([]byte, error) {
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

func (b *BoltBucket) Put(key string, data []byte) error {
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

func (b *BoltBucket) Delete(key string) error {
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

func (b *BoltBucket) CreateBucket(string) (Bucket, error) {
	return nil, errors.New("not implemented")
}
func (b *BoltBucket) GetBucket(string) (Bucket, error) { return nil, errors.New("not implemented") }
func (b *BoltBucket) DeleteBucket(string) error        { return errors.New("not implemented") }

// Test the interface
var (
	_ Bucket = &BoltBucket{}
)
