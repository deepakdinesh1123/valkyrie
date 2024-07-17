package cmd

import (
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/spf13/cobra"
)

var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start worker",
	Long:  `Start worker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logs.GetLogger()

		logger.Info().Msg("Starting worker")
		return nil
	},
}
