package cmd

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use: "add",

	RunE: func(cmd *cobra.Command, args []string) error {
		application.Logger.Info("HEELLO")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
