package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken        string `env:"BOT_TOKEN" env-required:"true"`
	IsTelegramDebug bool   `env:"TELEGRAM_DEBUG" env-default:"false"`
	IsAppDebug      bool   `env:"APP_DEBUG" env-default:"false"`
	DatabasePath    string `env:"DATABASE_PATH" env-required:"true"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't read env: %w", err)
	}

	return cfg, nil
}
