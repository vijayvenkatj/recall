package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Recall shell hooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := application.Install(); err != nil {
			return err
		}
		fmt.Println("✓ Installed Recall hook")
		fmt.Println()
		fmt.Println("Add the following line to your ~/.zshrc:")
		fmt.Println(`source "$HOME/.config/recall/hooks.zsh"`)
		fmt.Println()
		fmt.Println("Reload your shell:")
		fmt.Println(`source ~/.zshrc`)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
