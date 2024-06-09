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

type Bucketeer interface {
	CreateBucket(string) (Bucket, error)
	GetBucket(string) (Bucket, error)
	DeleteBucket(string) error
}

type Transaction interface {
	IsWritable() bool

	Bucketeer

	Rollback() error
	Commit() error
}

type KVPair interface {
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Delete(key string) error
}

type Bucket interface {
	Bucketeer
	KVPair
}
