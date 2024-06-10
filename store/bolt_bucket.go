package store

import (
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

func (b *BoltBucket) CreateBucket(name string) (Bucket, error) {
	bucket, err := b.bucket.CreateBucket([]byte(name))
	if err != nil {
		return nil, err
	}
	return newBoltBucket(bucket), nil
}

func (b *BoltBucket) GetBucket(name string) (Bucket, error) {
	bucket := b.bucket.Bucket([]byte(name))
	if bucket == nil {
		return nil, ErrBucketNotFound
	}
	return newBoltBucket(bucket), nil
}

func (b *BoltBucket) DeleteBucket(name string) error {
	return b.bucket.DeleteBucket([]byte(name))
}

// Test the interface
var (
	_ Bucket = &BoltBucket{}
)
