package app

import (
	"github.com/vijayvenkatj/recall/internal/config"
	"github.com/vijayvenkatj/recall/internal/repository"
	"go.uber.org/zap"
)

type App struct {
	Config config.Config
	Store  repository.Store
	Logger *zap.Logger
}

func New(config config.Config, store repository.Store, logger *zap.Logger) *App {
	return &App{
		Config: config,
		Store:  store,
		Logger: logger,
	}
}
