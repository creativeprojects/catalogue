package database

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/creativeprojects/catalogue/store"
)

type testStoreData struct {
	name  string
	store store.Store
}

func TestDatabase(t *testing.T) {
	// Prepare stores
	testStores := make([]testStoreData, 1)
	testStores[0] = testStoreData{"InMemory", store.NewMemoryStore()}

	// Add bolt store if the database path is set in the environment
	testPath := os.Getenv("DB_TEST_PATH")
	if testPath != "" {
		database := path.Join(testPath, "database_test.db")
		boltStore, err := store.NewBoltStore(database)
		if err == nil {
			defer os.Remove(database)
			testStores = append(testStores, testStoreData{"BoltDB", boltStore})
		} else {
			t.Logf("BoltDB tests cannot run: %s", err)
		}
	}

	for _, testData := range testStores {
		testData := testData // capture range variable
		database := NewDatabase(testData.store)
		t.Run(testData.name, func(t *testing.T) {
			// Don't run this suite in parallel, or the store will be closed immediately

			t.Run("TestInitAndStats", func(t *testing.T) {
				t.Parallel()

				database.Init()
				stats := database.Stats()
				assert.NotNil(t, stats.Created)
				assert.NotNil(t, stats.LastSaved)
				assert.WithinDuration(t, time.Now(), stats.Created, 10*time.Second)
			})
		})
		testData.store.Close()
	}
}
