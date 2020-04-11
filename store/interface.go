package store

type Store interface {
	GetBucket(string) (Bucket, error)

	GetKey(bucket, key string) ([]byte, error)
	GetKeyUint64(bucket, key string) (uint64, error)

	SetKey(bucket, key string, data []byte) error
	SetKeyUint64(bucket, key string, data uint64) error

	// Begin a transaction
	Begin(writable bool) (Transaction, error)
}

type Bucket interface {
	GetKey(key string) ([]byte, error)
	GetKeyUint64(key string) (uint64, error)

	SetKey(key string, data []byte) error
	SetKeyUint64(key string, data uint64) error
}

type Transaction interface {
	Rollback()
	Commit()
}
