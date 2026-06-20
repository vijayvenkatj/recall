package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/vijayvenkatj/recall/internal/config"
	"github.com/vijayvenkatj/recall/internal/llm"
	"github.com/vijayvenkatj/recall/internal/repository"
	"go.uber.org/zap"
)

type App struct {
	Config      config.Config
	Store       repository.Store
	Logger      *zap.Logger
	LLMProvider llm.Provider
}

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8E8E8E")).
			Italic(true)

	CommandListStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#8E8E8E"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#7D56F4")).
			PaddingLeft(2)
)

func New(config config.Config, store repository.Store, logger *zap.Logger) *App {
	llmClient, _ := llm.NewClient(config.LLMProvider, config.LLMAPIKey, config.LLMModel, config.LLMEndpoint)
	return &App{
		Config:      config,
		Store:       store,
		Logger:      logger,
		LLMProvider: llmClient,
	}
}
