// Entrypoint for CLI
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "development"

func NewRootCmd(v *viper.Viper) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "astrocyte",
		Short:   "matrix 2.0 server",
		Long:    "astrocyte - matrix 2.0 server",
		Version: version, // overriden by ldflags at build time
	}

	rootCmd.PersistentFlags().String("log-level", "info", "set log level")
	v.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	return rootCmd
}

func Execute() error {
	cfg, err := NewConfig()
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
