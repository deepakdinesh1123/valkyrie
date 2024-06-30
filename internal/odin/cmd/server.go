package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := server.NewServer()

		logs.Logger.Info().Msg("Listening on port 8080")
		err := http.ListenAndServe(":8080", server)
		if err != nil {
			logs.Logger.Err(err).Msg("Failed to start server")
			return err
		}
		return nil
	},
}
