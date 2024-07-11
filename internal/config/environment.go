package config

import (
	"github.com/spf13/viper"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

// Environment struct represents the configuration settings for the application.
type Environment struct {
	RABBITMQ_HOST string `mapstructure:"RABBITMQ_HOST"` // represents the RabbitMQ server host.
	RABBITMQ_PORT string `mapstructure:"RABBITMQ_PORT"` // represents the RabbitMQ server port.

	POSTGRES_HOST     string `mapstructure:"POSTGRES_HOST"`     // represents the PostgreSQL server host.
	POSTGRES_PORT     string `mapstructure:"POSTGRES_PORT"`     // represents the PostgreSQL server port.
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`     // represents the username for connecting to PostgreSQL.
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"` // represents the password for connecting to PostgreSQL.
	POSTGRES_DB       string `mapstructure:"POSTGRES_DB"`       // represents the name of the PostgreSQL database.
	POSTGRES_SSL_MODE string `mapstructure:"POSTGRES_SSL_MODE"` // represents the SSL mode for connecting to PostgreSQL.
	DB_MIGRATE        bool   `mapstructure:"DB_MIGRATE"`        // Indicates whether the database migrations should be run on startup
}

// EnvConfig holds the configuration settings for the application.

func GetEnvConfig() (*Environment, error) {
	logger := logs.GetLogger()
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	// Read configuration from file
	err := viper.ReadInConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to read configuration file")
		return nil, err
	}

	// Unmarshal configuration into EnvConfig struct
	var EnvConfig Environment
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	return &EnvConfig, nil
}
