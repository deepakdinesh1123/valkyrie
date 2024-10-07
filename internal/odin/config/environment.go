package config

import (
	"fmt"
	"log"
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

	ODIN_CONTAINER_ENGINE    string `mapstructure:"ODIN_CONTAINER_ENGINE"`    // represents the container engine used to execute code.
	ODIN_WORKER_EXECUTOR     string `mapstructure:"ODIN_WORKER_EXECUTOR"`     // represents the worker provider.
	ODIN_WORKER_CONCURRENCY  int32  `mapstructure:"ODIN_WORKER_CONCURRENCY"`  // represents the concurrency level for the worker.
	ODIN_WORKER_BUFFER_SIZE  int    `mapstructure:"ODIN_WORKER_BUFFER_SIZE"`  // represents the buffer size for the worker.
	ODIN_WORKER_TASK_TIMEOUT int    `mapstructure:"ODIN_WORKER_TASK_TIMEOUT"` // represents the task timeout.
	ODIN_WORKER_POLL_FREQ    int    `mapstructure:"ODIN_WORKER_POLL_FREQ"`    // represents the polling frequency for the worker in seconds.
	ODIN_WORKER_RUNTIME      string `mapstructure:"ODIN_WORKER_RUNTIME"`      // represents the default runtime for the worker containers (e.g. runc, crun).
	ODIN_WORKER_PODMAN_IMAGE string `mapstructure:"ODIN_WORKER_PODMAN_IMAGE"` // represents the default image for the podman worker containers.
	ODIN_WORKER_DOCKER_IMAGE string `mapstructure:"ODIN_WORKER_DOCKER_IMAGE"` // represents the default image for the docker worker containers.

	ODIN_LOG_LEVEL string `mapstructure:"ODIN_LOG_LEVEL"`

	ODIN_INFO_DIR           string
	ODIN_WORKER_DIR         string
	ODIN_WORKER_INFO_FILE   string
	ODIN_ENABLE_TELEMETRY   bool   `mapstructure:"ODIN_ENABLE_TELEMETRY"` // represents whether to enable OpenTelemetry for the server.
	ODIN_OTLP_ENDPOINT      string `mapstructure:"ODIN_OTLP_ENDPOINT"`    // represents the OpenTelemetry collector endpoint.
	ODIN_OTEL_RESOURCE_NAME string `mapstructure:"ODIN_OTEL_RESOURCE_NAME"`
	ODIN_EXPORT_LOGS        string `mapstructure:"ODIN_EXPORT_LOGS"`
	ODIN_ENVIRONMENT        string `mapstructure:"ODIN_ENVIRONMENT"` // represents the environment for the server (e.g. dev, staging, prod).

	ODIN_NIX_STORE string `mapstructure:"ODIN_NIX_STORE"` // represents the Nix store directory.

	ODIN_JOB_PRUNE_FREQ int `mapstructure:"ODIN_JOB_PRUNE_FREQ"` // represents the job prune frequency in hours.

	ODIN_SYSTEM_EXECUTOR_BASE_DIR string `mapstructure:"ODIN_SYSTEM_PROVIDER_BASE_DIR"` // represents the base directory for the system provider.
	ODIN_SYSTEM_EXECUTOR_CLEAN_UP bool   `mapstructure:"ODIN_SYSTEM_PROVIDER_CLEAN_UP"` // represents whether to clean up direcories created by the system provider.

	ODIN_USER_TOKEN  string `mapstructure:"ODIN_USER_TOKEN"`  // represents the secret key for the server.
	ODIN_ADMIN_TOKEN string `mapstructure:"ODIN_ADMIN_TOKEN"` // represents the admin token for the server.

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

	viper.SetDefault("ODIN_CONTAINER_ENGINE", "docker")
	viper.SetDefault("ODIN_WORKER_EXECUTOR", "system")
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

	viper.SetDefault("ODIN_SYSTEM_EXECUTOR_BASE_DIR", filepath.Join(os.TempDir(), "valkyrie"))
	viper.SetDefault("ODIN_SYSTEM_EXECUTOR_CLEAN_UP", true)

	viper.SetDefault("ODIN_JOB_PRUNE_FREQ", 1)

	viper.SetDefault("ODIN_LOG_LEVEL", "info")

	// Read configuration from file
	err := viper.ReadInConfig()
	if err != nil {
		log.Default().Println(".env file not found proceeding with defaults")
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
