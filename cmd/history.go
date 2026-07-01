package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

var numHistory int

var historyCmd = &cobra.Command{
	Use:   "history [query]",
	Short: "Browse and search your captured command history by session",
	Long: `Browse your captured sessions, or search for the ones containing a command.

  recall history            # recent sessions
  recall history docker     # sessions containing a command matching "docker"

Press enter to view a session's commands, / to filter, q to quit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		// Sync first so history reflects commands you just ran.
		if err := application.Sync(ctx); err != nil {
			return err
		}
		return application.History(ctx, strings.Join(args, " "), numHistory)
	},
}

func init() {
	historyCmd.Flags().IntVarP(&numHistory, "limit", "n", 20, "Number of sessions to show")
	rootCmd.AddCommand(historyCmd)
}
