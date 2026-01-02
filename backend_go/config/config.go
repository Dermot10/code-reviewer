package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedisAddr string `env:"REDIS_ADDR,required"`
}

func LoadConfig() (*Config, error) {
	props := Config{}
	if err := env.Parse(&props); err != nil {
		return nil, fmt.Errorf("error loading properties: %w", err)
	}
	return &props, nil
}
