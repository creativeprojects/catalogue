package database

import (
	"encoding/binary"
	"time"

	"github.com/creativeprojects/catalogue/store"
	"github.com/creativeprojects/catalogue/volume"
	"github.com/google/uuid"
)

const (
	BucketVolumes       = "catalogue-volumes"
	BucketStats         = "catalogue-stats"
	KeyDatabaseID       = "catalogue-id"
	KeyVersion          = "database-version"
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
	Version          Version
	TotalVolumes     uint64
	TotalDirectories uint64
	TotalFiles       uint64
	Created          time.Time
	LastSaved        time.Time
}

type Version struct {
	Major uint8
	Minor uint8
}

var (
	// CurrentVersion is the accepted database version
	CurrentVersion = Version{1, 0}
)

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
				stats.Put(KeyDatabaseID, bID)
			}
		}
		stats.Put(KeyVersion, versionToBytes(CurrentVersion))
		stats.Put(KeyTotalVolumes, uint64ToBytes(0))
		stats.Put(KeyTotalDirectories, uint64ToBytes(0))
		stats.Put(KeyTotalFiles, uint64ToBytes(0))

		now := time.Now()
		if bNow, err := now.MarshalBinary(); err == nil {
			stats.Put(KeyCreated, bNow)
			stats.Put(KeyLastSaved, bNow)
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
		stats.DatabaseID = bytesToUUID(bucket.Get(KeyDatabaseID))
		stats.Version = bytesToVersion(bucket.Get(KeyVersion))

		stats.TotalVolumes = bytesToUint64(bucket.Get(KeyTotalVolumes))
		stats.TotalDirectories = bytesToUint64(bucket.Get(KeyTotalDirectories))
		stats.TotalFiles = bytesToUint64(bucket.Get(KeyTotalFiles))

		stats.Created = bytesToTime(bucket.Get(KeyCreated))
		stats.LastSaved = bytesToTime(bucket.Get(KeyLastSaved))

		return nil
	})
	return stats
}

func (d *Database) IndexVolume(vol *volume.Volume) {
	// d.storage.Update(func(transaction store.Transaction) error {
	// 	bucket, err := transaction.GetBucket(BucketVolumes)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
}

func uint64ToBytes(value uint64) []byte {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return buffer
}

func int64ToBytes(value int64) []byte {
	return uint64ToBytes(uint64(value))
}

func versionToBytes(version Version) []byte {
	output := make([]byte, 2)
	output[0] = byte(version.Major)
	output[1] = byte(version.Minor)
	return output
}

func bytesToUint64(data []byte, err error) uint64 {
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(data)
}

func bytesToInt64(data []byte, err error) int64 {
	return int64(bytesToUint64(data, err))
}

func bytesToTime(data []byte, err error) time.Time {
	output := time.Time{}
	if err == nil {
		output.UnmarshalBinary(data)
	}
	return output
}

func bytesToUUID(data []byte, err error) uuid.UUID {
	output := uuid.UUID{}
	if err == nil {
		output.UnmarshalBinary(data)
	}
	return output
}

func bytesToVersion(data []byte, err error) Version {
	output := Version{1, 0}
	if err == nil {
		output.Major = uint8(data[0])
		output.Minor = uint8(data[1])
	}
	return output
}
