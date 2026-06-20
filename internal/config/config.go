package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver    string `mapstructure:"db_driver"`
	DBString    string `mapstructure:"db_string"`
	LogLevel    string `mapstructure:"log_level"`
	LogPath     string `mapstructure:"event_log_path"`
	LLMProvider string `mapstructure:"llm_provider"`
	LLMAPIKey   string `mapstructure:"llm_api_key"`
	LLMModel    string `mapstructure:"llm_model"`
	LLMEndpoint string `mapstructure:"llm_endpoint"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "recall")
	dataDir := filepath.Join(home, ".local", "share", "recall")
	defaultLogPath := filepath.Join(dataDir, "events.log")
	defaultDBString := filepath.Join(dataDir, "recall.db")

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return Config{}, err
	}

	viper.AutomaticEnv()

	_ = viper.BindEnv("db_driver", "DB_DRIVER")
	_ = viper.BindEnv("db_string", "DB_STRING")
	_ = viper.BindEnv("log_level", "LOG_LEVEL")
	_ = viper.BindEnv("event_log_path", "EVENT_LOG_PATH")
	_ = viper.BindEnv("llm_provider", "LLM_PROVIDER")
	_ = viper.BindEnv("llm_api_key", "LLM_API_KEY")
	_ = viper.BindEnv("llm_model", "LLM_MODEL")
	_ = viper.BindEnv("llm_endpoint", "LLM_ENDPOINT")

	// Load configuration file
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	if config.DBDriver == "" {
		config.DBDriver = "sqlite"
	}

	if config.DBString == "" {
		config.DBString = defaultDBString
	}

	if config.LogPath == "" {
		config.LogPath = defaultLogPath
	}

	// Resolve tildes in paths
	config.DBString = expandTilde(config.DBString)
	config.LogPath = expandTilde(config.LogPath)

	return config, nil
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
