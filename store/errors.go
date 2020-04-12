package store

import "errors"

// Errors
var (
	ErrTransactionReadonly = errors.New("Transaction is read-only")
	ErrNullPointerBucket   = errors.New("Null pointer bucket")
	ErrBucketNotFound      = errors.New("Bucket not found")
	ErrBucketNoName        = errors.New("Cannot use a blank name for a bucket")
	ErrBucketReadOnly      = errors.New("Cannot write or save a key, the bucket is read-only")
	ErrBucketNameExists    = errors.New("A bucket with this name already exists")
	ErrKeyNoName           = errors.New("Cannot use a blank name for a key")
	ErrKeyNotFound         = errors.New("Key not found")
)
