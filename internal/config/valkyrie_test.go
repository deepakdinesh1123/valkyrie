package config_test

import (
	"os"
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
)

func TestGetValkyrieConfig(t *testing.T) {
	_, err := config.GetValkyrieConfig()
	if err != nil {
		t.Errorf("Failed to get Valkyrie config: %v", err)
	}
}

func TestCreateValkyrieConfig(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Failed to get user home directory: %v", err)
	}
	configPath := home + "/.valkyrie/valkyrie.toml"
	os.Rename(configPath, configPath+".bak")
	_, err = config.CreateValkyrieConfig()
	if err != nil {
		t.Errorf("Failed to create Valkyrie config: %v", err)
	}
	os.Remove(configPath)
	os.Rename(configPath+".bak", configPath)
}
