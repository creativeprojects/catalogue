package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Volumes management",
	Long:  "List, add or delete volumes",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
