package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "odin",
	Short: "ODIN",
	Long:  `ODIN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Execute odin --help for more information on using odin.")
		return nil
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
