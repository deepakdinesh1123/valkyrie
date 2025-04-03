package store

import (
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var ripDBPath string
var valkyrieStoreConfig string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate valkyrie store",
	Long:  "Generate valkyrie store",
	RunE:  generatePackages,
}

func init() {
	GenerateCmd.Flags().StringVarP(&ripDBPath, "rip-db", "r", "", "The rippkgs db to use")
	GenerateCmd.Flags().StringVarP(&valkyrieStoreConfig, "config", "c", "", "valkyrie store config")
}

func generatePackages(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return err
	}
	logLevel := cmd.Flag("log-level").Value.String()
	config := logs.NewLogConfig(
		logs.WithLevel(logLevel),
		logs.WithExport(envConfig.EXPORT_LOGS),
		logs.WithSource("cli"),
	)
	logger := logs.GetLogger(config)
	err = store.GeneratePackages(ctx, valkyrieStoreConfig, ripDBPath, envConfig, logger)
	return err
}
