package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/cmd/server"
	"github.com/deepakdinesh1123/valkyrie/internal/cmd/store"
	"github.com/deepakdinesh1123/valkyrie/internal/cmd/worker"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
)

var RootCmd = &cobra.Command{
	Use:   "valkyrie",
	Short: "VALKYRIE",
	Long:  `VALKYRIE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Usage()
		return nil
	},
}

func Execute() {
	RootCmd.Execute()
}

func init() {
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		panic(err)
	}

	RootCmd.PersistentFlags().String("log-level", envConfig.LOG_LEVEL, "Log level")

	RootCmd.AddCommand(server.ServerCmd)
	RootCmd.AddCommand(worker.WorkerCmd)
	RootCmd.AddCommand(StandaloneCmd)
	RootCmd.AddCommand(store.StoreCmd)
	RootCmd.AddCommand(CompletionCmd)

	createDirs(envConfig)
}

func createDirs(envConfig *config.EnvConfig) error {
	dirs := []string{envConfig.INFO_DIR}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Printf("Failed to create directory %s: %v", dir, err)
				return err
			}
		}
	}
	return nil
}
