package app

import (
	"fmt"

	"github.com/caarlos0/env"
)

type (
	// Config - collects hole application configuration, separated by groups (example: ServerConfig - server settings)
	Config struct {
		Database PostgresConfig
		Logger   LoggerConfig
		Server   ServerConfig
	}
	// LoggerConfig - logger group settings
	LoggerConfig struct {
		Level string `env:"LOGGER_LEVEL" envDefault:"info"`
	}
	// ServerConfig - server group settings
	ServerConfig struct {
		Host            string `env:"SERVER_HOST"             envDefault:"0.0.0.0"`
		APIPort         int    `env:"SERVER_API_PORT"         envDefault:"8000"`
		HealthcheckPort int    `env:"SERVER_HEALTHCHECK_PORT" envDefault:"8081"`
		MetricsPort     int    `env:"SERVER_METRICS_PORT"     envDefault:"9090"`
	}
	// PostgresConfig - postgres database group settings
	PostgresConfig struct {
		Host     string `env:"POSTGRES_HOST"`
		Port     int    `env:"POSTGRES_PORT"`
		Database string `env:"POSTGRES_DATABASE"`
		Username string `env:"POSTGRES_USERNAME"`
		Password string `env:"POSTGRES_PASSWORD"`
	}
)

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Database,
		c.Username,
		c.Password,
	)
}

func NewConfig() (Config, error) {
	var (
		postgresCfg PostgresConfig
		loggerCfg   LoggerConfig
		serverCfg   ServerConfig
	)

	if err := env.Parse(&postgresCfg); err != nil {
		return Config{}, err
	}

	if err := env.Parse(&serverCfg); err != nil {
		return Config{}, err
	}

	if err := env.Parse(&loggerCfg); err != nil {
		return Config{}, err
	}

	return Config{
		Server:   serverCfg,
		Logger:   loggerCfg,
		Database: postgresCfg,
	}, nil
}
