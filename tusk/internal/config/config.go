package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	AppName string `mapstructure:"app_name"`
	Debug   bool   `mapstructure:"debug"`
}

var globalConfig *Config

// Load reads configuration from file, environment variables, and defaults
func Load(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("$HOME/.tusk")

	// Set defaults
	viper.SetDefault("app_name", "tusk")
	viper.SetDefault("debug", false)

	// Environment variables
	viper.SetEnvPrefix("TUSK")
	viper.AutomaticEnv()

	// If config path is provided, use it
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get returns the global configuration instance
func Get() *Config {
	if globalConfig == nil {
		return &Config{
			AppName: "tusk",
			Debug:   false,
		}
	}
	return globalConfig
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./configs/config.yaml"
	}
	return filepath.Join(home, ".tusk", "config.yaml")
}
