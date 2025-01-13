package store

import (
	"github.com/spf13/cobra"
)

var StoreCmd = &cobra.Command{
	Use:   "store",
	Short: "store command",
	Long:  `store command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	StoreCmd.AddCommand(GenerateCmd)
	StoreCmd.AddCommand(RealiseCmd)
}
