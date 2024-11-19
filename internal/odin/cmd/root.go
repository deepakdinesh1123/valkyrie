package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd/server"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd/store"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd/worker"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
)

var RootCmd = &cobra.Command{
	Use:   "odin",
	Short: "ODIN",
	Long:  `ODIN`,
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

	RootCmd.PersistentFlags().String("log-level", envConfig.ODIN_LOG_LEVEL, "Log level")

	RootCmd.AddCommand(server.ServerCmd)
	RootCmd.AddCommand(worker.WorkerCmd)
	RootCmd.AddCommand(StandaloneCmd)
	RootCmd.AddCommand(store.StoreCmd)

	createDirs(envConfig)
}

func createDirs(envConfig *config.EnvConfig) error {
	dirs := []string{envConfig.ODIN_INFO_DIR}
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
