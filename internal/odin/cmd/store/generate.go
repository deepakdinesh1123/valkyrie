package store

import (
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var ripDBPath string
var odinStoreConfig string

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate odin store",
	Long:  "Generate Odin store",
	RunE:  generatePackages,
}

func init() {
	GenerateCmd.Flags().StringVarP(&ripDBPath, "rip-db", "r", "", "The rippkgs db to use")
	GenerateCmd.Flags().StringVarP(&odinStoreConfig, "config", "c", "", "Odin store config")
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
		logs.WithExport(envConfig.ODIN_EXPORT_LOGS),
		logs.WithSource("cli"),
	)
	logger := logs.GetLogger(config)
	err = store.GeneratePackages(ctx, odinStoreConfig, ripDBPath, envConfig, logger)
	return err
}
