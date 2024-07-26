package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/server"
)

var StandaloneCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Odin standalone mode",
	Long:  `Start Odin in standalone mode`,
	RunE:  standaloneExec,
}

func standaloneExec(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	logger := logs.GetLogger()
	logger.Info().Msg("Starting Odin in standalone mode")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	done := make(chan bool, 1)

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get environment config")
		return err
	}

	applyMigrations, err := cmd.Flags().GetBool("migrate")
	if err != nil {
		logger.Err(err).Msg("Failed to get migrate flag")
		return err
	}

	dbConn, queries, err := db.GetDBConnection(ctx, true, envConfig, applyMigrations, sigChan, done, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to get database connection")
		cancel()
		return err
	}
	srv := server.NewServer(ctx, envConfig, dbConn, queries, logger)

	logger.Info().Msg("Starting Odin server")
	addr := fmt.Sprintf("%s:%s", envConfig.ODIN_SERVER_HOST, envConfig.ODIN_SERVER_PORT)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Err(err).Msg("Failed to start server")
	}

	<-done
	httpServer.Shutdown(ctx)
	logger.Info().Msg("Odin server stopped")
	return nil
}

func init() {
	StandaloneCmd.Flags().Bool("migrate", true, "Migrate database")
}
