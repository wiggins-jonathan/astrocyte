package cmd

import (
	"astrocyte/server"

	"github.com/spf13/cobra"
)

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "Port on which to serve astrocyte.")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"server"},
	Short:   "Serve astrocyte",
	Long:    "Serve astrocyte",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ListenAndServe will deal with port errors
		port, _ := cmd.Flags().GetInt("port")

		s := server.NewServer(
			server.WithDebug(debug),
			server.WithPort(port),
		)

		if err := s.Serve(); err != nil {
			return err
		}

		return nil
	},
}
