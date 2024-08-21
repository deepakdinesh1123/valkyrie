package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a job",
	Long:  `Execute a job`,
	RunE:  executeExec,
}

func executeExec(cmd *cobra.Command, args []string) error {
	baseURL := cmd.Flag("base-url").Value.String()

	client, err := api.NewClient(baseURL)
	if err != nil {
		return err
	}

	var req api.ExecutionRequest

	language := cmd.Flag("language").Value.String()
	file := cmd.Flag("file_path").Value.String()
	// arguments := cmd.Flag("args").Value.String()
	code := cmd.Flag("code").Value.String()
	flake := cmd.Flag("flake").Value.String()
	// dir := cmd.Flag("dir").Value.String()

	if language != "" {
		req.Language = language
		req.Environment.Set = true
		req.Environment.Value.Type = "ExecutionEnvironmentSpec"

		if code != "" && file != "" {
			return fmt.Errorf("language must be specified with either code or file")
		} else if code != "" {
			req.Code = code
		} else if file != "" {
			file_content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			req.Code = string(file_content)
		}
	} else if flake != "" {
		req.Environment.Set = true
		req.Environment.Value.Type = "Flake"
		req.Environment.Value.Flake = api.Flake(flake)
	} else {
		return fmt.Errorf("environment must be specified with either language or flake not both")
	}
	res, err := client.Execute(cmd.Context(), &req)
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *api.ExecuteOK:
		fmt.Println("Execution ID: ", res.ExecutionId)
	case *api.ExecuteBadRequest:
		fmt.Println(res.Message)
	case *api.ExecuteInternalServerError:
		fmt.Println(res.Message)
	}
	return nil
}

func init() {
	executeCmd.Flags().String("language", "", "Language to execute in")
	executeCmd.Flags().String("file_path", "", "File to execute")
	executeCmd.Flags().String("args", "", "Arguments to pass to the script")
	executeCmd.Flags().String("code", "", "Code to execute")
	executeCmd.Flags().String("flake", "", "Flake for the environment")
	executeCmd.Flags().String("dir", ".", "Path to the directory that contains the flake and script")
}
