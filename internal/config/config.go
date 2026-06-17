package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBString string `mapstructure:"DB_STRING"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	LogPath  string `mapstructure:"EVENT_LOG_PATH"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	home, _ := os.UserHomeDir()
	defaultLogPath := filepath.Join(home, ".local", "share", "recall", "events.log")

	viper.AutomaticEnv()

	_ = viper.BindEnv("DB_DRIVER")
	_ = viper.BindEnv("DB_STRING")
	_ = viper.BindEnv("LOG_LEVEL")
	_ = viper.BindEnv("EVENT_LOG_PATH")

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	if config.LogPath == "" {
		config.LogPath = defaultLogPath
	}

	return config, nil
}
