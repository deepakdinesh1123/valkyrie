package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type OdinServer struct {
	queries          *db.Queries
	envConfig        *config.EnvConfig
	executionService *execution.ExecutionService
	logger           *zerolog.Logger
	server           *api.Server
	connPool         *pgxpool.Pool
	tp               trace.TracerProvider
	mp               metric.MeterProvider
	otelShutdown     func(context.Context) error
}

func NewServer(ctx context.Context, envConfig *config.EnvConfig, standalone bool, applyMigrations bool, logger *zerolog.Logger) (*OdinServer, error) {
	connPool, queries, err := db.GetDBConnection(ctx, standalone, envConfig, applyMigrations, false, logger)
	if err != nil {
		return nil, err
	}
	executionService := execution.NewExecutionService(queries, envConfig, logger)
	otelShutdown, tp, mp, err := telemetry.SetupOTelSDK(ctx, "Odin Server", envConfig)
	if err != nil {
		logger.Err(err).Msg("Failed to setup OpenTelemetry")
		return nil, err
	}
	odinServer := &OdinServer{
		queries:          queries,
		executionService: executionService,
		envConfig:        envConfig,
		logger:           logger,
		connPool:         connPool,
		tp:               tp,
		mp:               mp,
		otelShutdown:     otelShutdown,
	}
	srv, err := api.NewServer(odinServer)
	if err != nil {
		return nil, err
	}

	odinServer.server = srv
	return odinServer, nil
}
