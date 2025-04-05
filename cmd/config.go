package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ConfigOption func(*viper.Viper) error

// NewConfig creates config with adjustable defaults
func NewConfig(options ...ConfigOption) (*viper.Viper, error) {
	v := viper.New()

	// environment variables
	configName := "astrocyte"
	v.SetEnvPrefix(strings.ToUpper(configName))
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // replace - with _
	v.AutomaticEnv()

	// config file
	v.SetConfigName(configName)

	configHome, exists := os.LookupEnv("XDG_CONFIG_HOME")
	if !exists {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving user home directory: %w", err)
		}

		configHome = filepath.Join(homeDir, ".config")
	}

	v.AddConfigPath(configHome)
	v.AddConfigPath(".")

	for _, option := range options {
		if err := option(v); err != nil {
			return nil, err
		}
	}

	if err := v.ReadInConfig(); err != nil { // read config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Error reading config file: %s: %w", configHome, err)
		}
	}

	return v, nil
}

// WithConfig allows the default file path to be overwritten with a custom path
func WithConfigPath(configPath string) ConfigOption {
	return func(v *viper.Viper) error {
		if configPath != "" {
			v.SetConfigFile(configPath)
		}
		return nil
	}
}
