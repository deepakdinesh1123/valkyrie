package cmd

import (
	"fmt"
	"strconv"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/spf13/cobra"
)

var (
	page   int32
	pageSz int32
)

var executionsCmd = &cobra.Command{
	Use:   "executions",
	Short: "Manage executions",
	Long:  `Manage executions`,
	RunE:  executionsExec,
}

func executionsExec(cmd *cobra.Command, args []string) error {
	_ = cmd.Usage()
	return nil
}

func init() {
	executionsCmd.AddCommand(executionsListCmd)
	executionsCmd.AddCommand(executionsResultsCmd)
	executionsCmd.AddCommand(execitionResultByIdCmd)
	executionsCmd.AddCommand(executionConfig)
	executionsCmd.AddCommand(deleteExecutionCmd)

	executionsCmd.PersistentFlags().Int32VarP(&page, "page", "p", 0, "Page number")
	executionsCmd.PersistentFlags().Int32VarP(&pageSz, "page-size", "s", 10, "Page size")
}

var executionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List executions",
	Long:  `List executions`,
	RunE:  executionsListExec,
}

func executionsListExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetAllExecutions(cmd.Context(), api.GetAllExecutionsParams{
		Page:     api.NewOptInt32(page),
		PageSize: api.NewOptInt32(pageSz),
	})
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.GetAllExecutionsOK:
		fmt.Println(res.Executions)
	case *api.GetAllExecutionsBadRequest:
		fmt.Println(res.Message)
	}
	return nil
}

var executionsResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "List execution results",
	Long:  `List execution results`,
	RunE:  executionsResultsExec,
}

func executionsResultsExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetAllExecutionResults(cmd.Context(), api.GetAllExecutionResultsParams{
		Page:     api.NewOptInt32(page),
		PageSize: api.NewOptInt32(pageSz),
	})
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.GetAllExecutionResultsOK:
		fmt.Println(res.Executions)
	case *api.GetAllExecutionResultsBadRequest:
		fmt.Println(res.Message)
	}
	return nil
}

var execitionResultByIdCmd = &cobra.Command{
	Use:   "get",
	Short: "Get execution result by id",
	Long:  `Get execution result by id`,
	RunE:  execitionResultByIdExec,
}

func execitionResultByIdExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	if len(args) == 0 {
		return fmt.Errorf("id is required")
	}
	jobId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return err
	}
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetExecutionResultsById(cmd.Context(), api.GetExecutionResultsByIdParams{
		JobId:    jobId,
		Page:     api.NewOptInt32(page),
		PageSize: api.NewOptInt32(pageSz),
	})
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.GetExecutionResultsByIdOK:
		fmt.Println(res.Executions)
	case *api.GetExecutionResultsByIdBadRequest:
		fmt.Println(res.Message)
	}
	return nil
}

var executionConfig = &cobra.Command{
	Use:   "config",
	Short: "Get execution config",
	Long:  `Get execution config`,
	RunE:  executionConfigExec,
}

func executionConfigExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.GetExecutionConfig(cmd.Context())
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

var deleteExecutionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete execution",
	Long:  `Delete execution`,
	RunE:  deleteExecutionExec,
}

func deleteExecutionExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()
	if len(args) == 0 {
		return fmt.Errorf("id is required")
	}
	jobId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return err
	}
	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}
	res, err := client.DeleteJob(cmd.Context(), api.DeleteJobParams{
		JobId: jobId,
	})
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *api.DeleteJobOK:
		fmt.Println(res)
	case *api.DeleteJobBadRequest:
		fmt.Println(res.Message)
	}
	return nil
}
