package store

type Store interface {
	// Begin a transaction
	Begin(writable bool) (Transaction, error)

	// Close the store
	Close()

	// Update wraps a write transaction around the function
	Update(func(transaction Transaction) error) error
	// View is a read-only view of a the database
	View(func(transaction Transaction) error) error
}

type Transaction interface {
	IsWritable() bool

	CreateBucket(string) (Bucket, error)
	GetBucket(string) (Bucket, error)
	DeleteBucket(string) error

	Rollback() error
	Commit() error
}

type Bucket interface {
	GetKey(key string) ([]byte, error)
	GetKeyUint64(key string) (uint64, error)

	SetKey(key string, data []byte) error
	SetKeyUint64(key string, data uint64) error

	DeleteKey(key string) error
}
