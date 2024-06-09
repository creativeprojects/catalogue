package store

import "sync"

type MemoryTransaction struct {
	store          *MemoryStore
	readBuckets    map[string]*MemoryBucket
	writeBuckets   map[string]*MemoryBucket
	deletedBuckets map[string]bool
	mutex          *sync.Mutex
	writeable      bool
	close          func()
	closing        *sync.Mutex
	closed         bool
}

func newMemoryTransaction(store *MemoryStore, writeable bool, close func()) *MemoryTransaction {
	return &MemoryTransaction{
		store:          store,
		readBuckets:    store.buckets,
		writeBuckets:   make(map[string]*MemoryBucket, 0),
		deletedBuckets: make(map[string]bool, 0),
		mutex:          &sync.Mutex{},
		writeable:      writeable,
		close:          close,
		closing:        &sync.Mutex{},
	}
}

func (t *MemoryTransaction) IsWritable() bool {
	return t.writeable
}

func (t *MemoryTransaction) Rollback() error {
	t.closing.Lock()
	defer t.closing.Unlock()

	if t.closed {
		return nil
	}
	t.close()
	t.closed = true
	return nil
}

func (t *MemoryTransaction) Commit() error {
	t.closing.Lock()
	defer t.closing.Unlock()

	if t.closed {
		return nil
	}

	if t.writeable {
		// Deleted buckets
		for name, deleted := range t.deletedBuckets {
			if deleted {
				t.store.deleteBucket(name)
			}
		}

		// Updated buckets
		for name, updated := range t.writeBuckets {
			t.saveBucket(name, updated)
		}
	}
	t.close()
	t.closed = true
	return nil
}

// CreateBucket returns a new bucket. Returns an error if the name already exists
func (t *MemoryTransaction) CreateBucket(bucket string) (Bucket, error) {
	if !t.writeable {
		return nil, ErrTransactionReadonly
	}
	if bucket == "" {
		return nil, ErrBucketNoName
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Check for updated buckets
	if _, found := t.writeBuckets[bucket]; found {
		return nil, ErrBucketNameExists
	}

	// Check for untouched existing bucket
	if _, found := t.readBuckets[bucket]; found {
		return nil, ErrBucketNameExists
	}

	// Remove from deleted bucket, in case it's there
	if _, found := t.deletedBuckets[bucket]; found {
		delete(t.deletedBuckets, bucket)
	}

	b := newMemoryBucket()
	t.writeBuckets[bucket] = b
	return b, nil
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
		return b, nil
	}
	return nil, ErrBucketNotFound
}

// DeleteBucket removes the bucket from memory
func (t *MemoryTransaction) DeleteBucket(bucket string) error {
	if !t.writeable {
		return ErrTransactionReadonly
	}
	if bucket == "" {
		return ErrBucketNoName
	}

	t.deletedBuckets[bucket] = true
	if _, ok := t.writeBuckets[bucket]; ok {
		delete(t.writeBuckets, bucket)
	}
	return nil
}

func (t *MemoryTransaction) saveBucket(name string, b *MemoryBucket) {
	storeBucket := t.store.getBucket(name)

	// Updated keys
	for key, value := range b.writeKeys {
		storeBucket.Put(key, value)
	}
	// Deleted keys
	for key, deleted := range b.deletedKeys {
		if deleted {
			storeBucket.Delete(key)
		}
	}
}

var (
	_ Transaction = &MemoryTransaction{}
)
