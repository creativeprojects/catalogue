package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStoreData struct {
	name  string
	store Store
}

func TestStores(t *testing.T) {
	// Prepare stores
	testStores := make([]testStoreData, 1)
	testStores[0] = testStoreData{"InMemory", NewMemoryStore()}

	for _, testData := range testStores {
		testData := testData // capture range variable
		t.Run(testData.name, func(t *testing.T) {
			t.Parallel()

			t.Run("TestCanCreateReadonlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotNil(t, tx)
				assert.False(t, tx.IsWritable())
				tx.Rollback()
			})

			t.Run("TestCanCreateWriteTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				if err != nil {
					t.Fatal(err)
				}

				assert.NotNil(t, tx)
				assert.True(t, tx.IsWritable())
				tx.Rollback()
			})

			t.Run("TestLoadEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				_, err = tx.GetBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestLoadUnknownBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				_, err = tx.GetBucket("unknown-bucket")
				assert.Equal(t, ErrBucketNotFound, err)
			})

		})
		testData.store.Close()
	}
}
