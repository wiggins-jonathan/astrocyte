// Entrypoint for CLI + utility functions
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "development"
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:     "astrocyte",
	Short:   "matrix 2.0 server",
	Long:    "astrocyte - matrix 2.0 server",
	Version: version, // overriden by ldflags at build time
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			fmt.Println("Debug mode enabled!")
		}
	},
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.PersistentFlags().BoolVarP(
		&debug, "debug", "d", false, "Enable Debug mode",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
