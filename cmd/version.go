package cmd

import (
	"fmt"

	"github.com/creativeprojects/catalogue/constants"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + constants.Catalogue,
	Long:  "Prints out " + constants.Catalogue + " version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constants.Description, constants.Version)
	},
}
