package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Recall version",
	// Override the root PersistentPreRunE so version works without config/DB.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("recall %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
