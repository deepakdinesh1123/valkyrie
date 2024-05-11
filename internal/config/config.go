// Package config provides functionality to manage environment configuration settings.
// It reads environment variables using viper and populates the Environment struct accordingly.
package config

import (
	"github.com/spf13/viper"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

// Environment struct represents the configuration settings for the application.
type Environment struct {
	REDIS_HOST            string `mapstructure:"REDIS_HOST"`            // represents the Redis server host.
	REDIS_PORT            string `mapstructure:"REDIS_PORT"`            // represents the Redis server port.
	POSTGRES_HOST         string `mapstructure:"POSTGRES_HOST"`         // represents the PostgreSQL server host.
	POSTGRES_PORT         string `mapstructure:"POSTGRES_PORT"`         // represents the PostgreSQL server port.
	POSTGRES_USER         string `mapstructure:"POSTGRES_USER"`         // represents the username for connecting to PostgreSQL.
	POSTGRES_PASSWORD     string `mapstructure:"POSTGRES_PASSWORD"`     // represents the password for connecting to PostgreSQL.
	POSTGRES_DB           string `mapstructure:"POSTGRES_DB"`           // represents the name of the PostgreSQL database.
	EXECUTION_ENVIRONMENT string `mapstructure:"EXECUTION_ENVIRONMENT"` // Indicates whether the execution environment is docker or k8s
	CONTAINERS            int    `mapstructure:"CONTAINERS"`            // Number of nix containers that will be spun up to execute code when the environment is docker
}

// EnvConfig holds the configuration settings for the application.
var EnvConfig Environment

// init initializes the configuration settings by reading environment variables using viper.
func init() {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	// Read configuration from file
	err := viper.ReadInConfig()
	if err != nil {
		logs.Logger.Err(err)
	}

	// Unmarshal configuration into EnvConfig struct
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		logs.Logger.Err(err)
	}
}
