package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "odinc",
	Short: "ODIN client",
	Long:  `ODIN client`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Usage()
		return nil
	},
}

func Execute() {
	RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(executeCmd)
	RootCmd.AddCommand(executionsCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(executionWorkersCmd)

	RootCmd.PersistentFlags().String("base-url", "http://localhost:8080", "Base URL")
}
