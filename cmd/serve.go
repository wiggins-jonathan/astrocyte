package cmd

import (
	"fmt"
	"log/slog"
	"net/url"

	"astrocyte/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewServeCmd creates the serve cmd for the CLI, binds its flags, & starts the
// server
func NewServeCmd(v *viper.Viper, l *slog.Logger) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Short:   "Serve astrocyte",
		Long:    "Serve astrocyte",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// ListenAndServe will deal with port errors
			port := v.GetInt("port")

			baseURL, err := url.Parse(v.GetString("base-url"))
			if err != nil {
				return fmt.Errorf("invalid base-url: %w", err)
			}

			s := server.NewServer(
				server.WithLogger(l),
				server.WithPort(port),
				server.WithBaseURL(baseURL),
			)

			return s.Serve()
		},
	}

	serveCmd.Flags().IntP("port", "p", 8080, "Port on which to serve astrocyte")
	v.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	serveCmd.Flags().String("base-url", "matrix.org", "Sets the root URL for the server")
	v.BindPFlag("base-url", serveCmd.Flags().Lookup("base-url"))

	return serveCmd
}
