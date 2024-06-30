package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geri",
	Short: "GERI",
	Long:  `GERI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Execute geri --help for more information on using geri.")
		return nil
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
}
