package cmd

import (
	"log/slog"

	"astrocyte/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

			s := server.NewServer(
				server.WithLogger(l),
				server.WithPort(port),
			)

			return s.Serve()
		},
	}

	serveCmd.Flags().IntP("port", "p", 8080, "Port on which to serve astrocyte.")
	v.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	return serveCmd
}
