package cmd

import (
	"fmt"
	"os"

	"github.com/creativeprojects/catalogue/constants"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new database",
	Long:  "Initializes a new empty database. The command will fail if a database file already exists.",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(rootFlags.Database); err == nil || os.IsExist(err) {
			log.WithField("file", rootFlags.Database).Error("Cannot initialize new database: file already exists")
		}

		db, err := bolt.Open(rootFlags.Database, 0600, nil)
		if err != nil {
			log.WithError(err).Error("Cannot open database")
		}
		defer db.Close()

		// Create admin bucket
		db.Update(func(tx *bolt.Tx) error {
			var err error

			b, err := tx.CreateBucket([]byte(constants.BucketAdmin))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			err = b.Put([]byte(constants.KeyVolumes), []byte("0"))
			if err != nil {
				return fmt.Errorf("put key: %s", err)
			}
			err = b.Put([]byte(constants.KeyDirectories), []byte("0"))
			if err != nil {
				return fmt.Errorf("put key: %s", err)
			}
			err = b.Put([]byte(constants.KeyFiles), []byte("0"))
			if err != nil {
				return fmt.Errorf("put key: %s", err)
			}
			return nil
		})
	},
}
