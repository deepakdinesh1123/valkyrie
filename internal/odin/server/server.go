package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
)

type OdinServer struct {
	queries          *db.Queries
	envConfig        *config.EnvConfig
	executionService *execution.ExecutionService
	logger           *zerolog.Logger
	server           *api.Server
}

func NewServer(ctx context.Context, envConfig *config.EnvConfig, standalone bool, applyMigrations bool, logger *zerolog.Logger) (*OdinServer, error) {
	queries, err := db.GetDBConnection(ctx, standalone, envConfig, applyMigrations, false, logger)
	if err != nil {
		return nil, err
	}
	executionService := execution.NewExecutionService(queries, envConfig, logger)
	odinServer := &OdinServer{
		queries:          queries,
		executionService: executionService,
		envConfig:        envConfig,
		logger:           logger,
	}
	srv, err := api.NewServer(odinServer)
	if err != nil {
		return nil, err
	}

	odinServer.server = srv
	return odinServer, nil
}
