package store

import "errors"

// Errors
var (
	ErrBucketNoName = errors.New("Cannot use a blank name for a bucket")
	ErrKeyNoName    = errors.New("Cannot use a blank name for a key")
	ErrKeyNotFound  = errors.New("Key not found")
)
