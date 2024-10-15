package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"

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
	ODIN_HOT_CONTAINER       int    `mapstructure:"ODIN_HOT_CONTAINER"`       // represents the buffer size for the worker.
	ODIN_WORKER_TASK_TIMEOUT int    `mapstructure:"ODIN_WORKER_TASK_TIMEOUT"` // represents the task timeout.
	ODIN_WORKER_POLL_FREQ    int    `mapstructure:"ODIN_WORKER_POLL_FREQ"`    // represents the polling frequency for the worker in seconds.
	ODIN_WORKER_RUNTIME      string `mapstructure:"ODIN_WORKER_RUNTIME"`      // represents the default runtime for the worker containers (e.g. runc, crun).
	ODIN_WORKER_PODMAN_IMAGE string `mapstructure:"ODIN_WORKER_PODMAN_IMAGE"` // represents the default image for the podman worker containers.
	ODIN_WORKER_DOCKER_IMAGE string `mapstructure:"ODIN_WORKER_DOCKER_IMAGE"` // represents the default image for the docker worker containers.

	ODIN_WORKER_MEMORY_LIMIT int64  `mapstructure:"ODIN_WORKER_MEMORY_LIMIT"`
	ODIN_WORKER_CPU_LIMIT    string `mapstructure:"ODIN_WORKER_CPU_LIMIT"`

	ODIN_LOG_LEVEL string `mapstructure:"ODIN_LOG_LEVEL"`

	ODIN_INFO_DIR           string
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

var envConfig *EnvConfig
var getEnvConfigOnce sync.Once

func GetEnvConfig() (*EnvConfig, error) {
	getEnvConfigOnce.Do(func() {
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

		viper.SetDefault("ODIN_CONTAINER_ENGINE", "podman")
		viper.SetDefault("ODIN_WORKER_EXECUTOR", "system")
		viper.SetDefault("ODIN_WORKER_CONCURRENCY", 10)
		viper.SetDefault("ODIN_HOT_CONTAINER", 5)
		viper.SetDefault("ODIN_WORKER_TASK_TIMEOUT", 120)
		viper.SetDefault("ODIN_WORKER_POLL_FREQ", 1)
		viper.SetDefault("ODIN_WORKER_RUNTIME", "runc")

		viper.SetDefault("ODIN_WORKER_MEMORY_LIMIT", 500)

		viper.SetDefault("ODIN_ENABLE_TELEMETRY", false)
		viper.SetDefault("ODIN_OTLP_ENDPOINT", "localhost:4317")
		viper.SetDefault("ODIN_OTEL_RESOURCE_NAME", "Odin")
		viper.SetDefault("ODIN_ENVIRONMENT", "dev")
		viper.SetDefault("ODIN_EXPORT_LOGS", "console")

		viper.SetDefault("ODIN_SYSTEM_EXECUTOR_BASE_DIR", filepath.Join(os.TempDir(), "valkyrie"))
		viper.SetDefault("ODIN_SYSTEM_EXECUTOR_CLEAN_UP", true)

		viper.SetDefault("ODIN_JOB_PRUNE_FREQ", 1)

		viper.SetDefault("ODIN_LOG_LEVEL", "info")

		viper.SetDefault("ODIN_WORKER_DOCKER_IMAGE", "odin")
		viper.SetDefault("ODIN_WORKER_PODMAN_IMAGE", "odin")

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
		_ = viper.ReadInConfig()

		// Unmarshal configuration into EnvConfig struct
		err := viper.Unmarshal(&envConfig)
		if err != nil {
			return
		}

		usr, err := user.Current()
		if err != nil {
			return
		}

		if usr != nil {
			envConfig.USER_HOME_DIR = filepath.Join("/home", usr.Name)
		} else {
			home_dir, err := os.UserHomeDir()
			if err != nil {

			}
			envConfig.USER_HOME_DIR = home_dir
		}
		envConfig.ODIN_INFO_DIR = filepath.Join(envConfig.USER_HOME_DIR, ".odin_info")
		envConfig.ODIN_WORKER_INFO_FILE = filepath.Join(envConfig.ODIN_INFO_DIR, "worker.json")
	})
	if envConfig != nil {
		return envConfig, nil
	}
	return nil, fmt.Errorf(" could not get env config")
}
