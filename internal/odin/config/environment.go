package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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

	ODIN_LOG_LEVEL string `mapstructure:"ODIN_LOG_LEVEL"`

	ODIN_INFO_DIR         string
	ODIN_WORKER_DIR       string
	ODIN_WORKER_INFO_FILE string

	ODIN_SYSTEM_PROVIDER_BASE_DIR string `mapstructure:"ODIN_SYSTEM_PROVIDER_BASE_DIR"` // represents the base directory for the system provider.
	ODIN_SYSTEM_PROVIDER_CLEAN_UP bool   `mapstructure:"ODIN_SYSTEM_PROVIDER_CLEAN_UP"` // represents whether to clean up direcories created by the system provider.

	USER_HOME_DIR string
}

// EnvConfig holds the configuration settings for the application.

func GetEnvConfig() (*EnvConfig, error) {
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

	viper.SetDefault("ODIN_SYSTEM_PROVIDER_BASE_DIR", "/tmp/valkyrie")
	viper.SetDefault("ODIN_SYSTEM_PROVIDER_CLEAN_UP", true)

	viper.SetDefault("ODIN_LOG_LEVEL", "info")

	// Read configuration from file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Unmarshal configuration into EnvConfig struct
	var EnvConfig EnvConfig
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		return nil, err
	}

	EnvConfig.USER_HOME_DIR, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	EnvConfig.ODIN_INFO_DIR = fmt.Sprintf("%s/%s", EnvConfig.USER_HOME_DIR, ".odin")
	EnvConfig.ODIN_WORKER_DIR = fmt.Sprintf("%s/%s", EnvConfig.ODIN_INFO_DIR, "worker")
	EnvConfig.ODIN_WORKER_INFO_FILE = fmt.Sprintf("%s/%s", EnvConfig.ODIN_WORKER_DIR, "worker-info.json")

	return &EnvConfig, nil
}
