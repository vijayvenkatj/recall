package app

import (
	"os"
	"path/filepath"

	"github.com/vijayvenkatj/recall/internal/assets"
)

func (app *App) Install() error {
	home, err := os.UserHomeDir()
	if err != nil {
		app.Logger.Error("failed to determine home directory")
		return err
	}

	configDir := filepath.Join(home, ".config", "recall")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		app.Logger.Sugar().Errorw("failed to create config directory", "path", configDir, "error", err)
		return err
	}

	hookPath := filepath.Join(configDir, "hooks.zsh")
	if err := os.WriteFile(hookPath, []byte(assets.HooksZsh), 0644); err != nil {
		app.Logger.Sugar().Errorw("failed to write hook file", "path", hookPath, "error", err)
		return err
	}

	return nil
}