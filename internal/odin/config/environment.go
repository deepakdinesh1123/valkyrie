package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// EnvConfig represents the configuration settings for the application.
type EnvConfig struct {
	POSTGRES_HOST            string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT            uint32 `mapstructure:"POSTGRES_PORT"`
	POSTGRES_USER            string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD        string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB              string `mapstructure:"POSTGRES_DB"`
	POSTGRES_SSL_MODE        string `mapstructure:"POSTGRES_SSL_MODE"`
	POSTGRES_STANDALONE_PATH string `mapstructure:"POSTGRES_STANDALONE_PATH"`

	ODIN_SERVER_HOST string `mapstructure:"ODIN_SERVER_HOST"`
	ODIN_SERVER_PORT string `mapstructure:"ODIN_SERVER_PORT"`

	ODIN_CONTAINER_ENGINE       string `mapstructure:"ODIN_CONTAINER_ENGINE"`
	ODIN_WORKER_EXECUTOR        string `mapstructure:"ODIN_WORKER_EXECUTOR"`
	ODIN_WORKER_SYSTEM_EXECUTOR string `mapstructure:"ODIN_WORKER_SYSTEM_EXECUTOR"`
	ODIN_WORKER_CONCURRENCY     int32  `mapstructure:"ODIN_WORKER_CONCURRENCY"`
	ODIN_HOT_CONTAINER          int    `mapstructure:"ODIN_HOT_CONTAINER"`
	ODIN_WORKER_TASK_TIMEOUT    int    `mapstructure:"ODIN_WORKER_TASK_TIMEOUT"`
	ODIN_WORKER_POLL_FREQ       int    `mapstructure:"ODIN_WORKER_POLL_FREQ"`
	ODIN_WORKER_RUNTIME         string `mapstructure:"ODIN_WORKER_RUNTIME"`
	ODIN_WORKER_PODMAN_IMAGE    string `mapstructure:"ODIN_WORKER_PODMAN_IMAGE"`
	ODIN_WORKER_DOCKER_IMAGE    string `mapstructure:"ODIN_WORKER_DOCKER_IMAGE"`
	ODIN_MAX_RETRIES            int    `mapstructure:"ODIN_MAX_RETRIES"`

	ODIN_WORKER_CONTAINER_MEMORY_LIMIT int64 `mapstructure:"ODIN_WORKER_CONTAINER_MEMORY_LIMIT"`

	ODIN_MEMORY_LIMIT float64 `mapstructure:"ODIN_MEMORY_LIMIT"`
	ODIN_CPU_LIMIT    float64 `mapstructure:"ODIN_CPU_LIMIT"`

	ODIN_LOG_LEVEL string `mapstructure:"ODIN_LOG_LEVEL"`

	ODIN_INFO_DIR           string
	ODIN_WORKER_INFO_FILE   string
	ODIN_ENABLE_TELEMETRY   bool   `mapstructure:"ODIN_ENABLE_TELEMETRY"`
	ODIN_OTLP_ENDPOINT      string `mapstructure:"ODIN_OTLP_ENDPOINT"`
	ODIN_OTEL_RESOURCE_NAME string `mapstructure:"ODIN_OTEL_RESOURCE_NAME"`
	ODIN_EXPORT_LOGS        string `mapstructure:"ODIN_EXPORT_LOGS"`
	ODIN_ENVIRONMENT        string `mapstructure:"ODIN_ENVIRONMENT"`

	ODIN_NIX_STORE                string `mapstructure:"ODIN_NIX_STORE"`
	ODIN_NIX_USER_ENVIRONMENT     string `mapstructure:"ODIN_NIX_USER_ENVIRONMENT"`
	ODIN_NIX_CHANNELS_ENVIRONMENT string `mapstructure:"ODIN_NIX_CHANNELS_ENVIRONMENT"`

	ODIN_JOB_PRUNE_FREQ int `mapstructure:"ODIN_JOB_PRUNE_FREQ"`

	ODIN_SYSTEM_EXECUTOR_BASE_DIR string `mapstructure:"ODIN_SYSTEM_EXECUTOR_BASE_DIR"`
	ODIN_SYSTEM_EXECUTOR_CLEAN_UP bool   `mapstructure:"ODIN_SYSTEM_PROVIDER_CLEAN_UP"`

	ODIN_USER_TOKEN  string `mapstructure:"ODIN_USER_TOKEN"`
	ODIN_ADMIN_TOKEN string `mapstructure:"ODIN_ADMIN_TOKEN"`

	RIPPKGS_BASE_URL string `mapstructure:"RIPPKGS_BASE_URL"`

	ODIN_BASE_DIR string `mapstructure:"ODIN_BASE_DIR"`
}

var (
	envConfig        *EnvConfig
	getEnvConfigOnce sync.Once
)

// GetEnvConfig initializes and retrieves the configuration settings.
func GetEnvConfig() (*EnvConfig, error) {
	getEnvConfigOnce.Do(func() {
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.SetConfigName(".env")
		viper.AutomaticEnv()

		setDefaults()

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Warning: Failed to read .env file: %v. Using defaults and environment variables.", err)
		}

		if err := viper.Unmarshal(&envConfig); err != nil {
			log.Fatalf("Failed to unmarshal configuration: %v", err)
		}

		if envConfig.POSTGRES_STANDALONE_PATH == "" {
			envConfig.POSTGRES_STANDALONE_PATH = filepath.Join(envConfig.ODIN_BASE_DIR, ".zango", "stdb")
		}

		envConfig.ODIN_INFO_DIR = filepath.Join(envConfig.ODIN_BASE_DIR, ".odin_info")
		_, err := os.Stat(envConfig.ODIN_INFO_DIR)
		if os.IsNotExist(err) {
			err = os.MkdirAll(envConfig.ODIN_INFO_DIR, 0755)
			if err != nil {
				log.Fatalf("Failed to created odin ino dir %s : %s", envConfig.ODIN_BASE_DIR, err)
			}
		}

		envConfig.ODIN_WORKER_INFO_FILE = filepath.Join(envConfig.ODIN_INFO_DIR, "worker.json")
	})

	if envConfig != nil {
		return envConfig, nil
	}
	return nil, fmt.Errorf("could not get env config")
}

func setDefaults() {
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
	viper.SetDefault("ODIN_WORKER_SYSTEM_EXECUTOR", "native")
	viper.SetDefault("ODIN_WORKER_CONCURRENCY", 10)
	viper.SetDefault("ODIN_HOT_CONTAINER", 5)
	viper.SetDefault("ODIN_WORKER_TASK_TIMEOUT", 120)
	viper.SetDefault("ODIN_WORKER_POLL_FREQ", 30)
	viper.SetDefault("ODIN_WORKER_RUNTIME", "runc")
	viper.SetDefault("ODIN_MAX_RETRIES", 5)

	viper.SetDefault("ODIN_WORKER_CONTAINER_MEMORY_LIMIT", 500)

	viper.SetDefault("ODIN_MEMORY_LIMIT", 75)
	viper.SetDefault("ODIN_CPU_LIMIT", 75)

	viper.SetDefault("ODIN_NIX_STORE", "/nix")

	viper.SetDefault("ODIN_ENABLE_TELEMETRY", false)
	viper.SetDefault("ODIN_OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("ODIN_OTEL_RESOURCE_NAME", "Odin")
	viper.SetDefault("ODIN_ENVIRONMENT", "dev")
	viper.SetDefault("ODIN_EXPORT_LOGS", "console")
	viper.SetDefault("ODIN_SYSTEM_EXECUTOR_BASE_DIR", filepath.Join(os.TempDir(), "valkyrie"))
	viper.SetDefault("ODIN_SYSTEM_EXECUTOR_CLEAN_UP", true)
	viper.SetDefault("ODIN_JOB_PRUNE_FREQ", 1)
	viper.SetDefault("ODIN_LOG_LEVEL", "info")
	viper.SetDefault("ODIN_WORKER_DOCKER_IMAGE", "odin:alpine")
	viper.SetDefault("ODIN_WORKER_PODMAN_IMAGE", "odin:alpine")
	viper.SetDefault("RIPPKGS_BASE_URL", "https://valnix-stage-bucket.s3.us-east-1.amazonaws.com")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		viper.SetDefault("ODIN_BASE_DIR", os.TempDir())
	}
	if homeDir == "/root" {
		viper.SetDefault("ODIN_BASE_DIR", os.TempDir())
	} else {
		viper.SetDefault("ODIN_BASE_DIR", homeDir)
	}
}
