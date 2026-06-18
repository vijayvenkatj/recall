package cmd

import (
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Recall and setup environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		return application.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
