package cmd

import (
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/spf13/cobra"
)

var executionWorkersCmd = &cobra.Command{
	Use:   "workers",
	Short: "Manage execution workers",
	Long:  `Manage execution workers`,
	RunE:  executionWorkersExec,
}

func executionWorkersExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetExecutionWorkers(cmd.Context(), api.GetExecutionWorkersParams{
		Page:     api.NewOptInt32(page),
		PageSize: api.NewOptInt32(pageSz),
	})
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.GetExecutionWorkersOK:
		fmt.Println(res.Workers)
	}
	return nil
}

func init() {
	executionWorkersCmd.PersistentFlags().Int32VarP(&page, "page", "p", 0, "Page number")
	executionWorkersCmd.PersistentFlags().Int32VarP(&pageSz, "page-size", "s", 10, "Page size")
}
