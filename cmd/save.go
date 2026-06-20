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

		limit := numSessions
		if !cmd.Flags().Changed("sessions") && application.Config.SaveSessionsLimit > 0 {
			limit = application.Config.SaveSessionsLimit
		}

		return application.SaveMemory(ctx, limit, numCommands)
	},
}

func init() {
	saveCmd.Flags().IntVarP(&numSessions, "sessions", "s", 10, "Number of recent sessions to show")
	saveCmd.Flags().IntVarP(&numCommands, "commands", "c", 10, "Number of recent commands to show per session")
	rootCmd.AddCommand(saveCmd)
}
