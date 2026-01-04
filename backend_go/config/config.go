package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedisAddr      string `env:"REDIS_ADDR,required"`
	DatabaseURL    string `env:"DATABASE_URL,required"`
	JWTSecret      string `env:"JWT_SECRET,required"`
	DBMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	DBMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
}

func LoadConfig() (*Config, error) {
	props := Config{}
	if err := env.Parse(&props); err != nil {
		return nil, fmt.Errorf("error loading properties: %w", err)
	}
	return &props, nil
}
