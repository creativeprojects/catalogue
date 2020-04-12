package database

import (
	"encoding/binary"
	"time"

	"github.com/google/uuid"

	"github.com/creativeprojects/catalogue/store"
)

const (
	BucketVolumes       = "catalogue-volumes"
	BucketStats         = "catalogue-stats"
	KeyDatabaseID       = "catalogue-id"
	KeyTotalVolumes     = "total-volumes"
	KeyTotalDirectories = "total-directories"
	KeyTotalFiles       = "total-files"
	KeyCreated          = "created"
	KeyLastSaved        = "last-saved"
)

type Database struct {
	storage store.Store
}

type Stats struct {
	DatabaseID       uuid.UUID
	TotalVolumes     uint64
	TotalDirectories uint64
	TotalFiles       uint64
	Created          *time.Time
	LastSaved        *time.Time
}

func NewDatabase(s store.Store) *Database {
	return &Database{
		storage: s,
	}
}

// Init a blank database
func (d *Database) Init() {
	d.storage.Update(func(transaction store.Transaction) error {
		_, err := transaction.CreateBucket(BucketVolumes)
		if err != nil {
			return err
		}
		stats, err := transaction.CreateBucket(BucketStats)
		if err != nil {
			return err
		}
		if ID, err := uuid.NewRandom(); err == nil {
			if bID, err := ID.MarshalBinary(); err == nil {
				stats.SetKey(KeyDatabaseID, bID)
			}
		}
		stats.SetKey(KeyTotalVolumes, Uint64ToBytes(0))
		stats.SetKey(KeyTotalDirectories, Uint64ToBytes(0))
		stats.SetKey(KeyTotalFiles, Uint64ToBytes(0))

		now := time.Now()
		if bNow, err := now.MarshalBinary(); err == nil {
			stats.SetKey(KeyCreated, bNow)
			stats.SetKey(KeyLastSaved, bNow)
		}
		return nil
	})
}

func (d *Database) Stats() Stats {
	stats := Stats{}
	d.storage.View(func(transaction store.Transaction) error {
		bucket, err := transaction.GetBucket(BucketStats)
		if err != nil {
			return err
		}
		if ID, err := bucket.GetKey(KeyDatabaseID); err == nil {
			stats.DatabaseID.UnmarshalBinary(ID)
		}

		stats.TotalVolumes, _ = bucket.GetKeyUint64(KeyTotalVolumes)
		stats.TotalDirectories, _ = bucket.GetKeyUint64(KeyTotalDirectories)
		stats.TotalFiles, _ = bucket.GetKeyUint64(KeyTotalFiles)

		stats.Created = &time.Time{}
		if created, err := bucket.GetKey(KeyCreated); err == nil {
			stats.Created.UnmarshalBinary(created)
		}

		stats.LastSaved = &time.Time{}
		if lastSaved, err := bucket.GetKey(KeyLastSaved); err == nil {
			stats.LastSaved.UnmarshalBinary(lastSaved)
		}

		return nil
	})
	return stats
}

func Uint64ToBytes(value uint64) []byte {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return buffer
}

func Int64ToBytes(value int64) []byte {
	return Uint64ToBytes(uint64(value))
}
