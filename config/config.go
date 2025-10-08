// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	KeyPath string `mapstructure:"key_path"`
}

// Load initializes and loads configuration using Viper
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("key_path", "$HOME/.key/key.pem")

	// Set config name and paths
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.key")

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("KEY")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Check if config file not found
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Create default config file
			if err := createDefaultConfig(); err != nil {
				return nil, err
			}
			// Try reading again after creating
			if err := viper.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Unmarshal config into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// createDefaultConfig creates a default config file in $HOME/.key/config.toml
func createDefaultConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".key")
	configPath := filepath.Join(configDir, "config.toml")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Create default config content
	defaultConfig := `# Key Manager Configuration
key_path = "$HOME/.key/key.pem"
`

	// Write config file
	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
