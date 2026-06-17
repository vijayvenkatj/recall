package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync command history into the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return application.Sync(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
