package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/creativeprojects/catalogue/database"
	"github.com/creativeprojects/catalogue/store"
	"github.com/pterm/pterm"

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
			pterm.Error.Printf("Database %q not found\n", rootFlags.Database)
			return
		}

		store, err := store.NewBoltStore(rootFlags.Database)
		if err != nil {
			pterm.Error.Printf("Cannot open database: %v\n", err)
			return
		}
		defer store.Close()

		db := database.NewDatabase(store)
		stats := db.Stats()
		fmt.Println("")
		fmt.Printf("     Database file:  %s\n", rootFlags.Database)
		fmt.Printf("                ID:  %s\n", stats.DatabaseID.String())
		fmt.Printf("           Version:  %d.%d\n", stats.Version.Major, stats.Version.Minor)
		fmt.Printf("           Created:  %s\n", stats.Created.Format(time.DateTime))
		fmt.Printf("        Last saved:  %s\n", stats.LastSaved.Format(time.DateTime))
		fmt.Printf("     Total volumes:  %d\n", stats.TotalVolumes)
		fmt.Printf(" Total directories:  %d\n", stats.TotalDirectories)
		fmt.Printf("       Total files:  %d\n", stats.TotalFiles)
		fmt.Println("")
	},
}
