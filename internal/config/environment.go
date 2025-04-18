package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// EnvConfig represents the configuration settings for the application.
type EnvConfig struct {
	ENABLE_EXECUTION bool `mapstructure:"ENABLE_EXECUTION"`
	ENABLE_SANDBOX   bool `mapstructure:"ENABLE_SANDBOX"`

	POSTGRES_HOST            string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT            uint32 `mapstructure:"POSTGRES_PORT"`
	POSTGRES_USER            string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD        string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB              string `mapstructure:"POSTGRES_DB"`
	POSTGRES_SSL_MODE        string `mapstructure:"POSTGRES_SSL_MODE"`
	POSTGRES_STANDALONE_PATH string `mapstructure:"POSTGRES_STANDALONE_PATH"`

	SERVER_HOST string `mapstructure:"SERVER_HOST"`
	SERVER_PORT string `mapstructure:"SERVER_PORT"`

	STORE_URL            string `mapstructure:"STORE_URL"`
	STORE_IMAGE          string `mapstructure:"STORE_IMAGE"`
	STORE_CONTAINER      string
	SANDBOX_NIXPKGS_PATH string
	SANDBOX_NIXPKGS_REV  string `mapstructure:"SANDBOX_NIXPKGS_REV"`

	RUNTIME             string `mapstructure:"RUNTIME"`
	CONTAINER_RUNTIME   string `mapstructure:"CONTAINER_RUNTIME"`
	WORKER_CONCURRENCY  int32  `mapstructure:"WORKER_CONCURRENCY"`
	HOT_CONTAINER       int    `mapstructure:"HOT_CONTAINER"`
	WORKER_TASK_TIMEOUT int    `mapstructure:"WORKER_TASK_TIMEOUT"`
	WORKER_POLL_FREQ    int    `mapstructure:"WORKER_POLL_FREQ"`
	EXECUTION_IMAGE     string `mapstructure:"EXECUTION_IMAGE"`
	MAX_RETRIES         int    `mapstructure:"MAX_RETRIES"`

	WORKER_CONTAINER_MEMORY_LIMIT int64 `mapstructure:"WORKER_CONTAINER_MEMORY_LIMIT"`

	MEMORY_LIMIT float64 `mapstructure:"MEMORY_LIMIT"`
	CPU_LIMIT    float64 `mapstructure:"CPU_LIMIT"`

	LOG_LEVEL string `mapstructure:"LOG_LEVEL"`

	INFO_DIR           string
	WORKER_INFO_FILE   string
	ENABLE_TELEMETRY   bool   `mapstructure:"ENABLE_TELEMETRY"`
	OTLP_ENDPOINT      string `mapstructure:"OTLP_ENDPOINT"`
	OTEL_RESOURCE_NAME string `mapstructure:"OTEL_RESOURCE_NAME"`
	EXPORT_LOGS        string `mapstructure:"EXPORT_LOGS"`
	ENVIRONMENT        string `mapstructure:"ENVIRONMENT"`

	NIX_STORE                string `mapstructure:"NIX_STORE"`
	NIX_USER_ENVIRONMENT     string `mapstructure:"NIX_USER_ENVIRONMENT"`
	NIX_CHANNELS_ENVIRONMENT string `mapstructure:"NIX_CHANNELS_ENVIRONMENT"`

	JOB_PRUNE_FREQ int `mapstructure:"JOB_PRUNE_FREQ"`

	USER_TOKEN  string `mapstructure:"USER_TOKEN"`
	ADMIN_TOKEN string `mapstructure:"ADMIN_TOKEN"`

	RIPPKGS_BASE_URL string `mapstructure:"RIPPKGS_BASE_URL"`

	SANDBOX_IMAGE string `mapstructure:"SANDBOX_IMAGE"`
	BASE_DIR      string `mapstructure:"BASE_DIR"`

	PY_INDEX string `mapstructure:"PY_INDEX"`
}

var (
	envConfig        *EnvConfig
	getEnvConfigOnce sync.Once
)

func LoadEnvFile(filename string) ([]string, error) {
	// Open the .env file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a slice to store key-value pairs
	var envPairs []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		// Trim spaces
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Add the pair in the format "key=value"
		envPairs = append(envPairs, fmt.Sprintf("%s=%s", key, value))
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envPairs, nil
}

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
			log.Fatalf("Failed to unmarshal .env configuration: %v", err)
		}

		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" && envConfig.CONTAINER_RUNTIME != "runc" {
			log.Fatalf("The specified container runtime %s in not supported in %s", envConfig.CONTAINER_RUNTIME, runtime.GOOS)
		}

		if envConfig.POSTGRES_STANDALONE_PATH == "" {
			envConfig.POSTGRES_STANDALONE_PATH = filepath.Join(envConfig.BASE_DIR, ".zango", "stdb")
		}

		envConfig.INFO_DIR = filepath.Join(envConfig.BASE_DIR, ".valkyrie_info")
		_, err := os.Stat(envConfig.INFO_DIR)
		if os.IsNotExist(err) {
			err = os.MkdirAll(envConfig.INFO_DIR, 0755)
			if err != nil {
				log.Fatalf("Failed to created valkyrie ino dir %s : %s", envConfig.BASE_DIR, err)
			}
		}

		envConfig.WORKER_INFO_FILE = filepath.Join(envConfig.INFO_DIR, "worker.json")
		envConfig.SANDBOX_NIXPKGS_PATH = fmt.Sprintf("/var/cache/nixpkgs/NixOS-nixpkgs-%s", envConfig.SANDBOX_NIXPKGS_REV[:7])
	})

	if envConfig != nil {
		return envConfig, nil
	}
	return nil, fmt.Errorf("could not get env config")
}

func setDefaults() {
	viper.SetDefault("ENABLE_EXECUTION", true)
	viper.SetDefault("ENABLE_SANDBOX", true)

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "thors")
	viper.SetDefault("POSTGRES_PASSWORD", "thorkell")
	viper.SetDefault("POSTGRES_DB", "valkyrie")
	viper.SetDefault("POSTGRES_SSL_MODE", "disable")

	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8080")

	viper.SetDefault("STORE_URL", "http://valkyrie-store:5000")
	viper.SetDefault("STORE_IMAGE", "valkyrie_store:0.0.1")
	viper.SetDefault("STORE_CONTAINER", "valkyrie-store")
	viper.SetDefault("SANDBOX_NIXPKGS_REV", "b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c") // pragma: allowlist secret

	viper.SetDefault("RUNTIME", "docker")
	viper.SetDefault("WORKER_CONCURRENCY", 10)
	viper.SetDefault("HOT_CONTAINER", 1)
	viper.SetDefault("WORKER_TASK_TIMEOUT", 120)
	viper.SetDefault("WORKER_POLL_FREQ", 30)
	viper.SetDefault("EXECUTION_IMAGE", "valkyrie_execution:0.0.1-ubuntu")
	viper.SetDefault("MAX_RETRIES", 5)

	viper.SetDefault("WORKER_CONTAINER_MEMORY_LIMIT", 500)

	viper.SetDefault("MEMORY_LIMIT", 75)
	viper.SetDefault("CPU_LIMIT", 75)

	viper.SetDefault("ENABLE_TELEMETRY", false)
	viper.SetDefault("OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("OTEL_RESOURCE_NAME", "valkyrie")
	viper.SetDefault("ENVIRONMENT", "dev")

	viper.SetDefault("EXPORT_LOGS", "console")
	viper.SetDefault("JOB_PRUNE_FREQ", 1)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("RIPPKGS_BASE_URL", "https://valnix-stage-bucket.s3.us-east-1.amazonaws.com")

	viper.SetDefault("SANDBOX_IMAGE", "valkyrie_sandbox:0.0.1-ubuntu")
	viper.SetDefault("PY_INDEX", "http://valkyrie-devpi:3141")

	// containerRuntime := ""
	// switch runtime.GOOS {
	// case "darwin":
	// 	containerRuntime = "runc"
	// case "linux":
	// 	containerRuntime = "runsc"
	// }
	viper.SetDefault("CONTAINER_RUNTIME", "runc")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		viper.SetDefault("BASE_DIR", os.TempDir())
	}
	if homeDir == "/root" {
		viper.SetDefault("BASE_DIR", os.TempDir())
	} else {
		viper.SetDefault("BASE_DIR", homeDir)
	}
}
