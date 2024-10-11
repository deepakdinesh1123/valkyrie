package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	ODIN_WORKER_CONCURRENCY  int32  `mapstructure:"ODIN_WORKER_CONCURRENCY"`  // represents the concurrency level for the worker.
	ODIN_WORKER_BUFFER_SIZE  int    `mapstructure:"ODIN_WORKER_BUFFER_SIZE"`  // represents the buffer size for the worker.
	ODIN_WORKER_TASK_TIMEOUT int    `mapstructure:"ODIN_WORKER_TASK_TIMEOUT"` // represents the task timeout.
	ODIN_WORKER_POLL_FREQ    int    `mapstructure:"ODIN_WORKER_POLL_FREQ"`    // represents the polling frequency for the worker in seconds.
	ODIN_WORKER_RUNTIME      string `mapstructure:"ODIN_WORKER_RUNTIME"`      // represents the runtime for the worker in seconds.

	ODIN_LOG_LEVEL string `mapstructure:"ODIN_LOG_LEVEL"`

	ODIN_INFO_DIR           string
	ODIN_WORKER_DIR         string
	ODIN_WORKER_INFO_FILE   string
	ODIN_ENABLE_TELEMETRY   bool   `mapstructure:"ODIN_ENABLE_TELEMETRY"` // represents whether to enable OpenTelemetry for the server.
	ODIN_OTLP_ENDPOINT      string `mapstructure:"ODIN_OTLP_ENDPOINT"`    // represents the OpenTelemetry collector endpoint.
	ODIN_OTEL_RESOURCE_NAME string `mapstructure:"ODIN_OTEL_RESOURCE_NAME"`
	ODIN_EXPORT_LOGS        string `mapstructure:"ODIN_EXPORT_LOGS"`
	ODIN_ENVIRONMENT        string `mapstructure:"ODIN_ENVIRONMENT"` // represents the environment for the server (e.g. dev, staging, prod).

	ODIN_JOB_PRUNE_FREQ int // represents the job prune frequency in hours.

	ODIN_SYSTEM_PROVIDER_BASE_DIR string `mapstructure:"ODIN_SYSTEM_PROVIDER_BASE_DIR"` // represents the base directory for the system provider.
	ODIN_SYSTEM_PROVIDER_CLEAN_UP bool   `mapstructure:"ODIN_SYSTEM_PROVIDER_CLEAN_UP"` // represents whether to clean up direcories created by the system provider.

	ODIN_SECRET_KEY string `mapstructure:"ODIN_SECRET_KEY"` // represents the secret key for the server.
	ODIN_USER_NAME  string `mapstructure:"ODIN_USER_NAME"`  // represents the user name for the server.
	ODIN_USER_PASS  string `mapstructure:"ODIN_USER_PASS"`  // represents the user password for the server.

	USER_HOME_DIR string

	// Packages generation database
	NIXOS_VERSION      string `mapstructure:"NIXOS_VERSION"`      // represents the NixOS version to use.
	DATABASE_HOST      string `mapstructure:"DATABASE_HOST"`      // represents the database host.
	DATABASE_PORT      string `mapstructure:"DATABASE_PORT"`      // represents the database port.
	DATABASE_PASSWORD  string `mapstructure:"DATABASE_PASSWORD"`  // represents the database password.
	DATABASE_SSL_MODE  string `mapstructure:"DATABASE_SSL_MODE"`  // represents the SSL mode for connecting to PostgreSQL.
	DATABASE_CONTAINER string `mapstructure:"DATABASE_CONTAINER"` // represents the database container name.
	DATABASE_USER      string `mapstructure:"DATABASE_USER"`      // represents the database user.
	DATABASE_NAME      string `mapstructure:"DATABASE_NAME"`      // represents the database name.
	DUMP_PATH          string `mapstructure:"DUMP_PATH"`          // represents the path for dumps inside the container.
	LOCAL_DUMP_PATH    string `mapstructure:"LOCAL_DUMP_PATH"`    // represents the local path for dumps.
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
	viper.SetDefault("ODIN_WORKER_POLL_FREQ", 1)
	viper.SetDefault("ODIN_WORKER_RUNTIME", "runc")

	viper.SetDefault("ODIN_ENABLE_TELEMETRY", false)
	viper.SetDefault("ODIN_OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("ODIN_OTEL_RESOURCE_NAME", "Odin")
	viper.SetDefault("ODIN_ENVIRONMENT", "dev")
	viper.SetDefault("ODIN_EXPORT_LOGS", "console")

	viper.SetDefault("ODIN_SYSTEM_PROVIDER_BASE_DIR", filepath.Join(os.TempDir(), "valkyrie"))
	viper.SetDefault("ODIN_SYSTEM_PROVIDER_CLEAN_UP", true)

	viper.SetDefault("ODIN_USER_NAME", "admin")
	viper.SetDefault("ODIN_USER_PASS", "admin")

	viper.SetDefault("ODIN_JOB_PRUNE_FREQ", 1)

	viper.SetDefault("ODIN_LOG_LEVEL", "info")

	viper.SetDefault("NIXOS_VERSION", "24.05")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5433")
	viper.SetDefault("DATABASE_SSL_MODE", "disable")
	viper.SetDefault("DATABASE_CONTAINER", "nixos-packages-db")
	viper.SetDefault("DATABASE_USER", "thors")
	viper.SetDefault("DATABASE_NAME", "nixos_packages")
	viper.SetDefault("DUMP_PATH", "/dumps")
	viper.SetDefault("LOCAL_DUMP_PATH", "./dumps")

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
