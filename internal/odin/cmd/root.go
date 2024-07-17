package cmd

import (
	"github.com/spf13/cobra"
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
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(WorkerCmd)
	RootCmd.AddCommand(StandaloneCmd)
}
