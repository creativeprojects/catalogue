package store

type Store interface {
	// Begin a transaction
	Begin(writable bool) (Transaction, error)

	// Close the store
	Close()
}

type Transaction interface {
	IsWritable() bool

	GetBucket(string) (Bucket, error)
	DeleteBucket(string) error

	Rollback()
	Commit()
}

type Bucket interface {
	GetKey(key string) ([]byte, error)
	GetKeyUint64(key string) (uint64, error)

	SetKey(key string, data []byte) error
	SetKeyUint64(key string, data uint64) error

	DeleteKey(key string) error
}
