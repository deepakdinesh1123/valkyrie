package store

import (
	"fmt"
	"os/exec"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache odin store",
	Long:  "Cache Odin store",
	RunE:  cacheStore,
}

func cacheStore(cmd *cobra.Command, args []string) error {

	_, err := exec.Command("cached-nix-shell", "--version").Output()

	if err != nil {
		// If there's an error (e.g., Vim is not installed), print the error
		return fmt.Errorf("cached-nix-shell is not installed or not found in the PATH")
	}

	ctx := cmd.Context()
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return logError("Error fetching environment config", err)
	}

	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(false),
		db.IsStandalone(false),
		db.IsWorker(false),
	)

	config := logs.NewLogConfig()
	logger := logs.GetLogger(config)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		return logError("Error getting DB connection", err)
	}

	pkgs, err := queries.GetAllPackages(ctx)
	if err != nil {
		return fmt.Errorf("could not get packages: %s", err)
	}

	for _, pkg := range pkgs {
		var pkgName string
		if pkg.Pkgtype == "system" {
			pkgName = pkg.Name
		} else {
			pkgName = fmt.Sprintf("%s.%s", pkg.Language.String, pkg.Name)
		}
		_, err := exec.Command("cached-nix-shell", "-p", pkgName, "--run", "exit").Output()

		if err != nil {
			// If there's an error (e.g., Vim is not installed), print the error
			return fmt.Errorf("error caching package %s: %s", pkg.Name, err)
		}

	}
	return nil
}
