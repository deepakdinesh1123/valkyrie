package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/database"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	DB      *pgx.Conn
	Queries *db.Queries

	Queue *mq.MessageQueue

	ValkyrieConfig   *config.ValkyrieConfig
	executionService *execution.ExecutionService
}

func NewServer(ctx context.Context, logger *zerolog.Logger) *api.Server {
	err := createQueues(logger)
	if err != nil {
		logger.Err(err).Msg("Failed to create queues")
		panic(err)
	}
	valkyrieConfig, err := config.GetValkyrieConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get valkyrie config")
	}

	DB, queries, err := database.GetDBConnection(ctx, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to Postgres")
		panic(err)
	}

	executionService := execution.NewExecutionService(queries, valkyrieConfig)
	server := &Server{
		DB:               DB,
		Queries:          queries,
		ValkyrieConfig:   valkyrieConfig,
		executionService: executionService,
	}
	srv, err := api.NewServer(server)
	if err != nil {
		logger.Err(err).Msg("Failed to create server")
		panic(err)
	}
	return srv
}

func createQueues(logger *zerolog.Logger) error {
	_, err := mq.NewQueue("execute", true, false, false, false, nil)
	if err != nil {
		logger.Err(err).Msg("Failed to create queue: execute")
		return err
	}
	return nil
}
