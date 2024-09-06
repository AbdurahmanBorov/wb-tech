package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	NatsConfig   NatsConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
}

type ServerConfig struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
}

type NatsConfig struct {
	ClusterID string `env:"NATSCLUSTERID"`
	ClientID  string `env:"NATSCLIENTID"`
	NatsURL   string `env:"NATSURL"`
	NatsChan  string `env:"NATSCHANNEL"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return cfg, nil
}
