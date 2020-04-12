package cmd

import (
	"os"

	"github.com/creativeprojects/catalogue/database"
	"github.com/creativeprojects/catalogue/store"

	"github.com/apex/log"
	"github.com/spf13/cobra"
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
			return
		}

		store, err := store.NewBoltStore(rootFlags.Database)
		if err != nil {
			log.WithError(err).Error("Cannot open database")
			return
		}
		defer store.Close()

		db := database.NewDatabase(store)
		db.Init()
	},
}
