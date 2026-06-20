package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/vijayvenkatj/recall/internal/assets"
	"github.com/vijayvenkatj/recall/internal/db/migrations"
)

func (app *App) Install() error {
	// 1. Ensure Data & Config Directories exist
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "recall")
	dataDir := filepath.Join(home, ".local", "share", "recall")

	dirs := []string{configDir, dataDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 2. Run Migrations
	fmt.Println("Running database migrations...")
	if err := app.Migrate(context.Background()); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// 3. Install Shell Hooks
	fmt.Println("Installing shell hooks...")
	shellPath := os.Getenv("SHELL")
	shellType := "zsh"
	if strings.Contains(shellPath, "fish") {
		shellType = "fish"
	} else if strings.Contains(shellPath, "bash") {
		shellType = "bash"
	}

	var hookContent string
	var hookFile string
	var rcFile string
	var sourceCmd string

	switch shellType {
	case "fish":
		hookContent = assets.HooksFish
		hookFile = "hooks.fish"
		rcFile = "~/.config/fish/config.fish"
		sourceCmd = fmt.Sprintf("source %s", filepath.Join(configDir, "hooks.fish"))
	case "bash":
		hookContent = assets.HooksBash
		hookFile = "hooks.bash"
		rcFile = "~/.bashrc"
		sourceCmd = fmt.Sprintf("source %s", filepath.Join(configDir, "hooks.bash"))
	default:
		hookContent = assets.HooksZsh
		hookFile = "hooks.zsh"
		rcFile = "~/.zshrc"
		sourceCmd = fmt.Sprintf("source %s", filepath.Join(configDir, "hooks.zsh"))
	}

	hookPath := filepath.Join(configDir, hookFile)
	if err := os.WriteFile(hookPath, []byte(hookContent), 0644); err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	fmt.Println("\n" + TitleStyle.Render(" INSTALLATION SUCCESSFUL "))
	fmt.Printf("\n1. Database and logs are at: %s\n", dataDir)
	fmt.Printf("2. Shell hooks installed to: %s\n", hookPath)
	fmt.Println("\n" + SelectedStyle.Render("Final Step:") + " Add the following line to your " + SelectedStyle.Render(rcFile) + ":")
	fmt.Printf("\n   %s\n\n", sourceCmd)
	fmt.Printf("Then restart your terminal or run: source %s\n", rcFile)

	return nil
}

func (app *App) Migrate(ctx context.Context) error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	db := app.Store.DB
	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}
