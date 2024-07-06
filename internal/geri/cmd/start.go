package cmd

import (
	"github.com/deepakdinesh1123/valkyrie/internal/geri/consumer"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start geri",
	Long:  `Start geri`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logs.GetLogger()
		ctx := cmd.Context()
		gconsumer, err := consumer.NewConsumer(ctx, logger)
		if err != nil {
			return err
		}
		gconsumer.Start(ctx, logger)
		return nil
	},
}
