package cmd

import (
	"context"
	"database/sql"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vijayvenkatj/recall/internal/app"
	"github.com/vijayvenkatj/recall/internal/config"
	"github.com/vijayvenkatj/recall/internal/repository"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

var application *app.App

var rootCmd = &cobra.Command{
	Use:   "recall [query]",
	Short: "Recall CLI - Search and manage your terminal memories",
	Long: `Recall is a CLI tool to capture terminal sessions and search through them using FTS5.
To search your memories, simply provide a query as an argument:
  recall "my search term"`,

	// Allow unknown commands to be treated as search queries
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args: cobra.ArbitraryArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		return application.Search(context.Background(), strings.Join(args, " "))
	},

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
