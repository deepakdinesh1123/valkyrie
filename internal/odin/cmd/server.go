package cmd

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logs.GetLogger()
		ctx := cmd.Context()
		server := server.NewServer(ctx, logger)

		logger.Info().Msg("Listening on port 8080")
		err := http.ListenAndServe(":8080", server)
		if err != nil {
			logger.Err(err).Msg("Failed to start server")
			return err
		}
		return nil
	},
}
