package store

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testBucketData struct {
	name   string
	bucket Bucket
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getBuckets() []testBucketData {
	data := make([]testBucketData, 1)

	data[0] = testBucketData{"InMemory", newMemoryBucket(nil, true)}
	return data
}

func TestLoadingUnknownKeyFromBucket(t *testing.T) {
	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.bucket.GetKey("some-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestSetKeyAndGetKeyFromBucket(t *testing.T) {
	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			var err error
			var value1, value2 []byte

			value1 = []byte("test data")
			err = data.bucket.SetKey("test-key", value1)
			if err != nil {
				t.Fatal(err)
			}

			value2, err = data.bucket.GetKey("test-key")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, value1, value2)
		})
	}
}

func TestLoadingUnknownUint64ValueFromBucket(t *testing.T) {
	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.bucket.GetKeyUint64("some-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestSetKeyAndGetUint64ValueFromBucket(t *testing.T) {

	for _, data := range getBuckets() {
		t.Run(data.name, func(t *testing.T) {
			var err error
			var value1, value2 uint64

			value1 = rand.Uint64()
			err = data.bucket.SetKeyUint64("test-key", value1)
			if err != nil {
				t.Fatal(err)
			}

			value2, err = data.bucket.GetKeyUint64("test-key")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, value1, value2)
		})
	}
}
