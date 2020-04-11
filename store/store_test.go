package store

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStoreData struct {
	name  string
	store Store
}

func getStores() []testStoreData {
	data := make([]testStoreData, 1)

	data[0] = testStoreData{"InMemory", NewMemoryStore()}
	return data
}

func TestCannotGetEmptyBucketName(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.store.GetBucket("")
			assert.Equal(t, ErrBucketNoName, err)
		})
	}
}

func TestCannotGetEmptyKeyName(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.store.GetKey("bucket", "")
			assert.Equal(t, ErrKeyNoName, err)
		})
	}
}

func TestLoadingUnknownKeyFromStore(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.store.GetKey("my-bucket", "some-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestSetKeyAndGetKeyFromStore(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			var err error
			var value1, value2 []byte

			value1 = []byte("test data")
			err = data.store.SetKey("my-bucket", "test-key", value1)
			if err != nil {
				t.Fatal(err)
			}

			value2, err = data.store.GetKey("my-bucket", "test-key")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, value1, value2)
		})
	}
}

func TestSetKeyAndGetKeyFromDifferentBucket(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			var err error
			var value []byte

			value = []byte("test data")
			err = data.store.SetKey("my-first-bucket", "test-key", value)
			if err != nil {
				t.Fatal(err)
			}

			_, err = data.store.GetKey("my-second-bucket", "test-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestLoadingUnknownUint64ValueFromStore(t *testing.T) {
	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			_, err := data.store.GetKeyUint64("some-store", "some-key")
			assert.Equal(t, ErrKeyNotFound, err)
		})
	}
}

func TestSetKeyAndGetUint64ValueFromStore(t *testing.T) {

	for _, data := range getStores() {
		t.Run(data.name, func(t *testing.T) {
			var err error
			var value1, value2 uint64

			value1 = rand.Uint64()
			err = data.store.SetKeyUint64("my-store", "test-key", value1)
			if err != nil {
				t.Fatal(err)
			}

			value2, err = data.store.GetKeyUint64("my-store", "test-key")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, value1, value2)
		})
	}
}
