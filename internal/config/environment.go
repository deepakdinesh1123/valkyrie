package config

import (
	"github.com/spf13/viper"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

// Environment struct represents the configuration settings for the application.
type EnvConfig struct {
	POSTGRES_HOST     string `mapstructure:"POSTGRES_HOST"`     // represents the PostgreSQL server host.
	POSTGRES_PORT     uint32 `mapstructure:"POSTGRES_PORT"`     // represents the PostgreSQL server port.
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`     // represents the username for connecting to PostgreSQL.
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"` // represents the password for connecting to PostgreSQL.
	POSTGRES_DB       string `mapstructure:"POSTGRES_DB"`       // represents the name of the PostgreSQL database.
	POSTGRES_SSL_MODE string `mapstructure:"POSTGRES_SSL_MODE"` // represents the SSL mode for connecting to PostgreSQL.

	ODIN_SERVER_HOST string `mapstructure:"ODIN_SERVER_HOST"` // represents the host on which the Odin server will listen.
	ODIN_SERVER_PORT string `mapstructure:"ODIN_SERVER_PORT"` // represents the port on which the Odin server will listen.

	ODIN_WORKER_PROVIDER     string `mapstructure:"ODIN_WORKER_PROVIDER"`     // represents the worker provider.
	ODIN_WORKER_CONCURRENCY  int    `mapstructure:"ODIN_WORKER_CONCURRENCY"`  // represents the concurrency level for the worker.
	ODIN_WORKER_BUFFER_SIZE  int    `mapstructure:"ODIN_WORKER_BUFFER_SIZE"`  // represents the buffer size for the worker.
	ODIN_WORKER_TASK_TIMEOUT int    `mapstructure:"ODIN_WORKER_TASK_TIMEOUT"` // represents the task timeout.
	ODIN_WORKER_POLL_FREQ    int    `mapstructure:"ODIN_WORKER_POLL_FREQ"`    // represents the polling frequency for the worker in seconds.
}

// EnvConfig holds the configuration settings for the application.

func GetEnvConfig() (*EnvConfig, error) {
	logger := logs.GetLogger()
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "thors")
	viper.SetDefault("POSTGRES_PASSWORD", "thorkell")
	viper.SetDefault("POSTGRES_DB", "valkyrie")
	viper.SetDefault("POSTGRES_SSL_MODE", "disable")

	viper.SetDefault("ODIN_SERVER_HOST", "0.0.0.0")
	viper.SetDefault("ODIN_SERVER_PORT", "8080")

	viper.SetDefault("ODIN_WORKER_PROVIDER", "system")
	viper.SetDefault("ODIN_WORKER_CONCURRENCY", 10)
	viper.SetDefault("ODIN_WORKER_BUFFER_SIZE", 100)
	viper.SetDefault("ODIN_WORKER_TASK_TIMEOUT", 30)
	viper.SetDefault("ODIN_WORKER_POLL_FREQ", 5)

	// Read configuration from file
	err := viper.ReadInConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to read configuration file")
		return nil, err
	}

	// Unmarshal configuration into EnvConfig struct
	var EnvConfig EnvConfig
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	return &EnvConfig, nil
}
