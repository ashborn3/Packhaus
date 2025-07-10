package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PSWD,required"`
	DBName   string `env:"DB_NAME,required"`
}

func LoadConfig() (*Configuration, error) {
	cfg := Configuration{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading environment vars through .env: %s", err.Error())
	}

	return &cfg, nil
}
