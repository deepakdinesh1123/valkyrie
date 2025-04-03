package server

import (
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server command",
	Long:  `Server command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	ServerCmd.AddCommand(ServerStartCmd)
}
