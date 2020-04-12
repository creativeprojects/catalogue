package store

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	UnknownBucketKey = "unknown-bucket"
)

type testStoreData struct {
	name  string
	store Store
}

func TestClosingStoreShouldPanicIfTransactionIsRunning(t *testing.T) {
	store := NewMemoryStore()
	store.Begin(false)
	assert.Panics(t, store.Close)
}

func TestStores(t *testing.T) {
	// Prepare stores
	testStores := make([]testStoreData, 1)
	testStores[0] = testStoreData{"InMemory", NewMemoryStore()}

	// Add bolt store if the database path is set in the environment
	ramdisk := os.Getenv("RAMDISK")
	if ramdisk != "" {
		database := path.Join(ramdisk, "test.db")
		boltStore, err := NewBoltStore(database)
		if err == nil {
			defer os.Remove(database)
			testStores = append(testStores, testStoreData{"BoltDB", boltStore})
		} else {
			t.Logf("BoltDB tests cannot run: %s", err)
		}
	}

	for _, testData := range testStores {
		testData := testData // capture range variable
		t.Run(testData.name, func(t *testing.T) {
			// Don't run this suite in parallel, or the store will be closed immediately

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

			t.Run("TestCreateEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				_, err = tx.CreateBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestDeleteEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				err = tx.DeleteBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestLoadUnknownBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				_, err = tx.GetBucket(UnknownBucketKey)
				assert.Equal(t, ErrBucketNotFound, err)
			})

			t.Run("TestCannotCreateBucketInReadOnlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				_, err = tx.CreateBucket(UnknownBucketKey)
				assert.Equal(t, ErrTransactionReadonly, err)
			})

			t.Run("TestCannotDeleteBucketInReadOnlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				err = tx.DeleteBucket(UnknownBucketKey)
				assert.Equal(t, ErrTransactionReadonly, err)
			})

			t.Run("TestCannotSetKeysInReadonlyBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				// Isolate transaction
				func() {
					tx, err := testData.store.Begin(true)
					if err != nil {
						t.Fatal(err)
					}
					defer tx.Rollback()

					_, err = tx.CreateBucket(name)
					if err != nil {
						t.Fatal(err)
					}
					tx.Commit()
				}()

				tx, err := testData.store.Begin(false)
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()

				b, err := tx.GetBucket(name)
				if err != nil {
					t.Fatalf("Error getting bucket %s: %s", name, err)
				}
				err = b.SetKey("something", []byte("is something"))
				assert.Equal(t, ErrBucketReadOnly, err)
			})

		})
		testData.store.Close()
	}
}
