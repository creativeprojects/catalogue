package cmd

import (
	"fmt"
	"os"

	"github.com/creativeprojects/catalogue/constants"
	"github.com/spf13/cobra"
)

type RootFlags struct {
	Verbose  bool
	Database string
}

var (
	rootCmd = &cobra.Command{
		Use:   constants.Catalogue,
		Short: constants.Description,
		Long:  `An offline file catalogue with fast search`,
		Run: func(cmd *cobra.Command, args []string) {
			// set log level here
		},
	}

	rootFlags RootFlags
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Database, "database", "d", "catalogue.db", "database file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
