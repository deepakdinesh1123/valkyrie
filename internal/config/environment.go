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

	EXEC_ENV     string `mapstructure:"EXEC_ENV"`     // represents the execution environment for the application.
	CONCURRENCY  int    `mapstructure:"CONCURRENCY"`  // represents the concurrency level for the worker.
	BUFFER_SIZE  int    `mapstructure:"BUFFER_SIZE"`  // represents the buffer size for the worker.
	TASK_TIMEOUT int    `mapstructure:"TASK_TIMEOUT"` // represents the task timeout.
}

// EnvConfig holds the configuration settings for the application.

func GetEnvConfig() (*EnvConfig, error) {
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
	var EnvConfig EnvConfig
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	return &EnvConfig, nil
}
