package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type ValkyrieConfig struct {
	Provider string `toml:"environment"` // Provider
	Geri     Geri   `toml:"geri"`        // Geri configuration
}

type Geri struct {
	ExecutionEnvironment string `toml:"execution_environment"` // Execution environment
	Concurrency          int    `toml:"concurrency"`           // Number of concurrent workers
	BufferSize           int    `toml:"buffer_size"`           // Buffer size for the task queue
	TaskTimeout          int    `toml:"task_timeout"`          // Timeout in seconds
}

var (
	configFile = "valkyrie.toml"
)

// GetValkyrieConfig retrieves the Valkyrie configuration from the user's home directory.
//
// It first tries to locate the configuration file at ~/.valkyrie/valkyrie.toml.
// If the file does not exist, it calls CreateValkyrieConfig to create a new configuration file.
//
// Returns a ValkyrieConfig struct and an error if any error occurred during the process.
func GetValkyrieConfig() (*ValkyrieConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := home + "/.valkyrie/" + configFile

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config, err := CreateValkyrieConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	var config *ValkyrieConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return config, err
	}
	return config, nil
}

// CreateValkyrieConfig creates the Valkyrie configuration file in the user's home directory.
//
// Returns a pointer to the newly created ValkyrieConfig struct and an error if any error occurred during the process.
func CreateValkyrieConfig() (*ValkyrieConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(home + "/.valkyrie")
	if os.IsNotExist(err) {
		err = os.Mkdir(home+"/.valkyrie", 0755)
		if err != nil {
			return nil, err
		}
	}
	configPath := home + "/.valkyrie/" + configFile

	config := &ValkyrieConfig{
		Provider: "docker",
		Geri: Geri{
			ExecutionEnvironment: "NoSysbox",
			Concurrency:          1,
			BufferSize:           10,
			TaskTimeout:          300,
		},
	}
	configToml, err := toml.Marshal(config)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write([]byte(configToml))
	if err != nil {
		return nil, err
	}
	return config, nil
}
