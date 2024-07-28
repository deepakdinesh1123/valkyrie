package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
)

type Server struct {
	queries          *db.Queries
	envConfig        *config.EnvConfig
	executionService *execution.ExecutionService
}

func NewServer(ctx context.Context, envConfig *config.EnvConfig, queries *db.Queries, logger *zerolog.Logger) *api.Server {
	executionService := execution.NewExecutionService(queries, envConfig)
	server := &Server{
		queries:          queries,
		executionService: executionService,
		envConfig:        envConfig,
	}
	srv, err := api.NewServer(server)
	if err != nil {
		logger.Err(err).Msg("Failed to create server")
		panic(err)
	}
	return srv
}
