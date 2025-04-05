// Entrypoint for CLI
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "development"

// NewRootCmd creates the root command for the CLI & binds global &  persistent
// flags inherited by all subcommands
func NewRootCmd(v *viper.Viper) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "astrocyte",
		Short:   "matrix 2.0 server",
		Long:    "astrocyte - matrix 2.0 server",
		Version: version, // overriden by ldflags at build time
	}

	rootCmd.PersistentFlags().String("log-level", "info", "set log level")
	v.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	// even though we don't use viper/cobra to parse for the config file path
	// we still want to bind the flag to get correct usage().
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to the config file")
	v.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	return rootCmd
}

// Entrypoint for the CLI. Bootstraps config, logger, & registers subcommands
func Execute() error {
	configPath, err := getConfigFlag()
	if err != nil {
		return err
	}

	cfg, err := NewConfig(WithConfigPath(configPath))
	if err != nil {
		return fmt.Errorf("Error initializing config: %w", err)
	}

	logLevel := viper.GetString("log-level")
	logger := NewLogger(WithLogLevel(logLevel))

	rootCmd := NewRootCmd(cfg)
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	subCommands := []func(*viper.Viper, *slog.Logger) *cobra.Command{
		NewServeCmd,
	}

	for _, cmd := range subCommands {
		rootCmd.AddCommand(cmd(cfg, logger))
	}

	return rootCmd.Execute()
}

// getConfigFlag() parses CLI args in case the user wants to change the config
// file path. This must be done at the very beginning of startup to avoid a
// chicken/egg scenario. This spaghetti is why people use cobra/viper
func getConfigFlag() (string, error) {
	for i, arg := range os.Args {
		if arg == "--config" || arg == "-c" {
			if i+1 >= len(os.Args) { // user passed flag but no argument
				return "", fmt.Errorf("Error: flag needs an argument: %s", arg)
			}

			return os.Args[i+1], nil
		}

		// deal with an equals sign
		if strings.HasPrefix(arg, "--config=") || strings.HasPrefix(arg, "-c=") {
			filepath := strings.SplitN(arg, "=", 2)[1]
			filepath = strings.TrimSpace(filepath)

			if filepath == "" {
				return "", fmt.Errorf("Error: flag needs an argument: %s", arg)
			}

			return filepath, nil
		}
	}

	// deal with environment variable
	filepath, exists := os.LookupEnv("ASTROCYTE_CONFIG")
	if !exists {
		return "", nil // no flag or environment variable passed
	}

	return filepath, nil
}
