package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	DBPath      string `env:"DB_PATH" envDefault:"./data/database.sqlite"`
	LottoAPIKey string `env:"LOTTO_API_KEY"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			slog.Info("No .env file found, defaulting to environment variables")
		} else {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
