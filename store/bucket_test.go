package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBucketData struct {
	name   string
	bucket Bucket
}

func getBuckets() []testBucketData {
	data := make([]testBucketData, 1)

	data[0] = testBucketData{"InMemory", newMemoryBucket(nil, true)}
	return data
}

func TestLoadingUnknownKeyFromBucket(t *testing.T) {
	t.Parallel()

	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.bucket.Get("some-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestSetKeyAndGetKeyFromBucket(t *testing.T) {
	t.Parallel()

	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			t.Parallel()

			var err error
			var value1, value2 []byte

			value1 = []byte("test data")
			err = data.bucket.Put("test-key", value1)
			if err != nil {
				t.Fatal(err)
			}

			value2, err = data.bucket.Get("test-key")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, value1, value2)
		})
	}
}
