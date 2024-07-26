package worker

import (
	"github.com/spf13/cobra"
)

var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Worker command",
	Long:  `Worker command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	WorkerCmd.AddCommand(WorkerStartCmd)
}
