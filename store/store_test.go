package store

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	UnknownBucketKey = "unknown-bucket"
)

type testStoreData struct {
	name  string
	store Store
}

func TestClosingStoreShouldPanicIfTransactionIsRunning(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	store.Begin(false)
	assert.Panics(t, store.Close)
}

func TestStores(t *testing.T) {
	t.Parallel()

	// Prepare stores
	testStores := make([]testStoreData, 0, 2)
	testStores = append(testStores, testStoreData{"InMemory", NewMemoryStore()})

	// Add bolt store if the database path is set in the environment
	testPath := os.Getenv("DB_TEST_PATH")
	if testPath != "" {
		database := path.Join(testPath, "store_test.db")
		boltStore, err := NewBoltStore(database)
		if err == nil {
			t.Cleanup(func() {
				boltStore.Close()
				_ = os.Remove(database)
			})
			testStores = append(testStores, testStoreData{"BoltDB", boltStore})
		} else {
			t.Logf("BoltDB tests cannot run: %s", err)
		}
	}

	for _, testData := range testStores {
		t.Run(testData.name, func(t *testing.T) {

			t.Run("TestCanCreateReadonlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				require.NoError(t, err)

				assert.NotNil(t, tx)
				assert.False(t, tx.IsWritable())
				tx.Rollback()
			})

			t.Run("TestCanCreateTwoReadonlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx1, err := testData.store.Begin(false)
				require.NoError(t, err)
				assert.NotNil(t, tx1)

				tx2, err := testData.store.Begin(false)
				require.NoError(t, err)
				assert.NotNil(t, tx2)

				assert.False(t, tx1.IsWritable())
				assert.False(t, tx2.IsWritable())
				tx1.Rollback()
				tx2.Rollback()
			})

			t.Run("TestCanCreateWriteTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				require.NoError(t, err)

				assert.NotNil(t, tx)
				assert.True(t, tx.IsWritable())
				tx.Rollback()
			})

			t.Run("TestCanCreateReadWriteTransaction", func(t *testing.T) {
				// don't run this test in parallel (risk of deadlock)
				tx1, err := testData.store.Begin(false)
				require.NoError(t, err)
				assert.NotNil(t, tx1)

				tx2, err := testData.store.Begin(true)
				require.NoError(t, err)
				assert.NotNil(t, tx2)

				assert.False(t, tx1.IsWritable())
				assert.True(t, tx2.IsWritable())
				tx1.Rollback()
				tx2.Rollback()
			})

			t.Run("TestCanCreateWriteReadTransaction", func(t *testing.T) {
				// don't run this test in parallel (risk of deadlock)
				tx1, err := testData.store.Begin(true)
				require.NoError(t, err)
				assert.NotNil(t, tx1)

				tx2, err := testData.store.Begin(false)
				require.NoError(t, err)
				assert.NotNil(t, tx2)

				assert.True(t, tx1.IsWritable())
				assert.False(t, tx2.IsWritable())
				tx1.Rollback()
				tx2.Rollback()
			})

			t.Run("TestCannotLoadEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				require.NoError(t, err)
				defer tx.Rollback()

				_, err = tx.GetBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestCannotCreateEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				_, err = tx.CreateBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestCannotDeleteEmptyNameBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				err = tx.DeleteBucket("")
				assert.Equal(t, ErrBucketNoName, err)
			})

			t.Run("TestCannotLoadUnknownBucket", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				require.NoError(t, err)
				defer tx.Rollback()

				_, err = tx.GetBucket(UnknownBucketKey)
				assert.Equal(t, ErrBucketNotFound, err)
			})

			t.Run("TestCannotCreateBucketInReadOnlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				require.NoError(t, err)
				defer tx.Rollback()

				_, err = tx.CreateBucket(UnknownBucketKey)
				assert.Equal(t, ErrTransactionReadonly, err)
			})

			t.Run("TestCannotDeleteBucketInReadOnlyTransaction", func(t *testing.T) {
				t.Parallel()
				tx, err := testData.store.Begin(false)
				require.NoError(t, err)
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
					require.NoError(t, err)
					defer tx.Rollback()

					_, err = tx.CreateBucket(name)
					require.NoError(t, err)
					tx.Commit()
				}()

				tx, err := testData.store.Begin(false)
				require.NoError(t, err)
				defer tx.Rollback()

				b, err := tx.GetBucket(name)
				require.NoError(t, err)
				err = b.Put("something", []byte("is something"))
				assert.Equal(t, ErrBucketReadOnly, err)
			})

			t.Run("TestLoadingUnknownKeyFromBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				bucket, err := tx.CreateBucket(name)
				require.NoError(t, err)

				_, err = bucket.Get("some-key")
				assert.Equal(t, ErrKeyNotFound, err)
			})

			t.Run("TestSetKeyAndGetKeyFromBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				bucket, err := tx.CreateBucket(name)
				require.NoError(t, err)

				var value1, value2 []byte

				value1 = []byte("test data")
				err = bucket.Put("test-key", value1)
				require.NoError(t, err)

				value2, err = bucket.Get("test-key")
				require.NoError(t, err)
				assert.Equal(t, value1, value2)
			})

			t.Run("TestCreateBucketInBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				bucket, err := tx.CreateBucket(name)
				require.NoError(t, err)

				_, err = bucket.CreateBucket("sub-bucket")
				require.NoError(t, err)
			})

			t.Run("TestGetBucketInBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				bucket, err := tx.CreateBucket(name)
				require.NoError(t, err)

				_, err = bucket.CreateBucket("sub-bucket")
				require.NoError(t, err)

				_, err = bucket.GetBucket("sub-bucket")
				require.NoError(t, err)
			})

			t.Run("TestDeleteBucketInBucket", func(t *testing.T) {
				t.Parallel()
				name := path.Base(t.Name())

				tx, err := testData.store.Begin(true)
				require.NoError(t, err)
				defer tx.Rollback()

				bucket, err := tx.CreateBucket(name)
				require.NoError(t, err)

				_, err = bucket.CreateBucket("sub-bucket")
				require.NoError(t, err)

				err = bucket.DeleteBucket("sub-bucket")
				require.NoError(t, err)

				_, err = bucket.GetBucket("sub-bucket")
				require.ErrorIs(t, err, ErrBucketNotFound)
			})

		})
		testData.store.Close()
	}
}
