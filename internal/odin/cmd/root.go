package cmd

import (
	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd/server"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd/worker"
)

var RootCmd = &cobra.Command{
	Use:   "odin",
	Short: "ODIN",
	Long:  `ODIN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Usage()
		return nil
	},
}

func Execute() {
	RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(server.ServerCmd)
	RootCmd.AddCommand(worker.WorkerCmd)
	RootCmd.AddCommand(StandaloneCmd)
}
