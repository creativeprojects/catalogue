package cmd

import (
	"fmt"
	"os"

	"github.com/creativeprojects/catalogue/database"
	"github.com/creativeprojects/catalogue/store"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Database statistics",
	Long:  "Display some simple database statistics.",
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := os.Stat(rootFlags.Database); os.IsNotExist(err) {
			log.WithField("file", rootFlags.Database).Error("Database not found")
			return
		}

		store, err := store.NewBoltStore(rootFlags.Database)
		if err != nil {
			log.WithError(err).Error("Cannot open database")
			return
		}
		defer store.Close()

		db := database.NewDatabase(store)
		stats := db.Stats()
		fmt.Println("")
		fmt.Printf("          Database:  %s\n", rootFlags.Database)
		fmt.Printf("                ID:  %s\n", stats.DatabaseID.String())
		fmt.Printf("           Version:  %d.%d\n", stats.Version.Major, stats.Version.Minor)
		fmt.Printf("           Created:  %s\n", stats.Created)
		fmt.Printf("        Last saved:  %s\n", stats.LastSaved)
		fmt.Printf("     Total volumes:  %d\n", stats.TotalVolumes)
		fmt.Printf(" Total directories:  %d\n", stats.TotalDirectories)
		fmt.Printf("       Total files:  %d\n", stats.TotalFiles)
		fmt.Println("")
	},
}
