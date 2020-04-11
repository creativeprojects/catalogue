package store

import "sync"

type MemoryTransaction struct {
	store          *MemoryStore
	readBuckets    map[string]map[string][]byte
	writeBuckets   map[string]*MemoryBucket
	deletedBuckets map[string]bool
	mutex          *sync.Mutex
	writable       bool
	close          func()
	closing        *sync.Mutex
	closed         bool
}

func newMemoryTransaction(store *MemoryStore, writable bool, close func()) *MemoryTransaction {
	return &MemoryTransaction{
		store:          store,
		readBuckets:    store.buckets,
		writeBuckets:   make(map[string]*MemoryBucket, 0),
		deletedBuckets: make(map[string]bool, 0),
		mutex:          &sync.Mutex{},
		writable:       writable,
		close:          close,
		closing:        &sync.Mutex{},
	}
}

func (t *MemoryTransaction) IsWritable() bool {
	return t.writable
}

func (t *MemoryTransaction) Rollback() {
	t.closing.Lock()
	defer t.closing.Unlock()

	if t.closed {
		return
	}
	t.close()
	t.closed = true
}

func (t *MemoryTransaction) Commit() {
	t.closing.Lock()
	defer t.closing.Unlock()

	if t.closed {
		return
	}

	if t.writable {
		// Deleted buckets
		for name, deleted := range t.deletedBuckets {
			if deleted {
				delete(t.store.buckets, name)
			}
		}

		// Updated buckets
		for name, updated := range t.writeBuckets {
			t.saveBucket(name, updated)
		}
	}
	t.close()
	t.closed = true
}

// GetBucket returns a bucket from its name. If it does not exists, a new empty bucket will be returned.
func (t *MemoryTransaction) GetBucket(bucket string) (Bucket, error) {
	if bucket == "" {
		return nil, ErrBucketNoName
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Check for deleted bucket first
	if deleted, ok := t.deletedBuckets[bucket]; ok && deleted {
		return nil, ErrBucketNotFound
	}

	// Check for updated buckets
	if b, ok := t.writeBuckets[bucket]; ok {
		return b, nil
	}

	// Check for untouched existing bucket
	if b, ok := t.readBuckets[bucket]; ok {
		return NewMemoryBucket(b, t.writable), nil
	}
	return nil, ErrBucketNotFound
}

func (t *MemoryTransaction) DeleteBucket(bucket string) error {
	t.deletedBuckets[bucket] = true
	return nil
}

func (t *MemoryTransaction) saveBucket(name string, b *MemoryBucket) {
	//
}

var (
	_ Transaction = &MemoryTransaction{}
)
