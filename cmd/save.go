package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	numSessions int
	numCommands int
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Create a memory from a recent session",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		// Automatically sync before saving to ensure latest history is present
		if err := application.Sync(ctx); err != nil {
			return err
		}
		return application.SaveMemory(ctx, numSessions, numCommands)
	},
}

func init() {
	saveCmd.Flags().IntVarP(&numSessions, "sessions", "s", 3, "Number of recent sessions to show")
	saveCmd.Flags().IntVarP(&numCommands, "commands", "c", 10, "Number of recent commands to show per session")
	rootCmd.AddCommand(saveCmd)
}
