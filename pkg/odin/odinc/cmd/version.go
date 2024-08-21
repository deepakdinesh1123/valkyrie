package cmd

import (
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version",
	Long:  `Get version`,
	RunE:  versionExec,
}

func versionExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetVersion(cmd.Context())
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.GetVersionOK:
		fmt.Println(res.Version)
	}
	return nil
}
