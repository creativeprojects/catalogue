package store

import (
	bolt "go.etcd.io/bbolt"
)

type BoltTransaction struct {
	tx *bolt.Tx
}

func newBoltTransaction(tx *bolt.Tx) *BoltTransaction {
	return &BoltTransaction{
		tx: tx,
	}
}

func (t *BoltTransaction) IsWritable() bool {
	return t.tx.Writable()
}

func (t *BoltTransaction) Rollback() {
	t.tx.Rollback()
}

func (t *BoltTransaction) Commit() {
	t.tx.Commit()
}

// CreateBucket returns a new bucket. Returns an error if the name already exists
func (t *BoltTransaction) CreateBucket(bucket string) (Bucket, error) {
	if !t.tx.Writable() {
		return nil, ErrTransactionReadonly
	}
	if bucket == "" {
		return nil, ErrBucketNoName
	}

	b, err := t.tx.CreateBucket([]byte(bucket))
	if err != nil {
		return nil, err
	}
	return newBoltBucket(b), nil
}

// GetBucket returns a bucket from its name. If it does not exists, a new empty bucket will be returned.
func (t *BoltTransaction) GetBucket(bucket string) (Bucket, error) {
	if bucket == "" {
		return nil, ErrBucketNoName
	}

	b := t.tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, ErrBucketNotFound
	}
	return newBoltBucket(b), nil
}

// DeleteBucket removes the bucket from memory
func (t *BoltTransaction) DeleteBucket(bucket string) error {
	if !t.tx.Writable() {
		return ErrTransactionReadonly
	}
	if bucket == "" {
		return ErrBucketNoName
	}

	return t.tx.DeleteBucket([]byte(bucket))
}

var (
	_ Transaction = &BoltTransaction{}
)
