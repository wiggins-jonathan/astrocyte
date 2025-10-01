// Entrypoint for CLI
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "development"

type contextKey string

const loggerKey contextKey = "logger"

// NewRootCmd creates the root command for the CLI & binds global & persistent
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

	rootCmd.PersistentFlags().String("log-format", "text", "set log format")
	v.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))

	// even though we don't use viper/cobra to parse for the config file path
	// we still want to bind the flag to get correct usage().
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to the config file")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		logLevel, err := cmd.Flags().GetString("log-level")
		if err != nil {
			return fmt.Errorf("Failed to parse log-level: %w", err)
		}

		logFormat, err := cmd.Flags().GetString("log-format")
		if err != nil {
			return fmt.Errorf("Failed to parse log-format: %w", err)
		}

		v.Set("log-level", logLevel)
		v.Set("log-format", logFormat)

		logger, err := NewLogger(WithLevel(logLevel), WithFormat(logFormat))
		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), loggerKey, logger)
		cmd.SetContext(ctx)

		return nil
	}

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

	rootCmd := NewRootCmd(cfg)
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	subCommands := []func(*viper.Viper) *cobra.Command{
		NewServeCmd,
	}

	for _, cmd := range subCommands {
		rootCmd.AddCommand(cmd(cfg))
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

// getLogger uses context to pass the logger from the CLI to the log constructor
// & then to the server middleware
func getLogger(cmd *cobra.Command) *slog.Logger {
	if logger, ok := cmd.Context().Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	// Fallback: This should ideally only happen in tests or error cases
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
