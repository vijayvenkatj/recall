package cmd

import (
	"database/sql"
	"os"

	"github.com/spf13/cobra"
	"github.com/vijayvenkatj/recall/internal/app"
	"github.com/vijayvenkatj/recall/internal/config"
	"github.com/vijayvenkatj/recall/internal/repository"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

var application *app.App

var rootCmd = &cobra.Command{
	Use:   "recall",
	Short: "Recall CLI",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if application != nil {
			return nil
		}

		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		logger, err := zap.NewProduction()
		if err != nil {
			return err
		}

		db, err := sql.Open(config.DBDriver, config.DBString)
		if err != nil {
			return err
		}

		store := repository.New(db)
		application = app.New(config, *store, logger)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
