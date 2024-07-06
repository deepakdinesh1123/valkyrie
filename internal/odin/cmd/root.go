package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "odin",
	Short: "ODIN",
	Long:  `ODIN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Execute odin --help for more information on using odin.")
		return nil
	},
}

func Execute() {
	RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
