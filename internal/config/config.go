package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBString string `mapstructure:"DB_STRING"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	_ = viper.BindEnv("DB_DRIVER")
	_ = viper.BindEnv("DB_STRING")
	_ = viper.BindEnv("LOG_LEVEL")

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
