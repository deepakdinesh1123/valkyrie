package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type OdinServer struct {
	queries          db.Store
	envConfig        *config.EnvConfig
	executionService *execution.ExecutionService
	logger           *zerolog.Logger
	server           *api.Server
	tp               trace.TracerProvider
	mp               metric.MeterProvider
	otelShutdown     func(context.Context) error
}

// NewServer creates a new OdinServer instance with the provided configuration.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - envConfig: The configuration of the OdinServer.
// - standalone: A boolean indicating whether the server is standalone.
// - applyMigrations: A boolean indicating whether to apply migrations.
// - logger: The logger for the OdinServer.
//
// Returns:
// - *OdinServer: The newly created OdinServer instance.
// - error: An error if the OdinServer could not be created.
func NewServer(ctx context.Context, envConfig *config.EnvConfig, standalone bool, applyMigrations bool, logger *zerolog.Logger) (*OdinServer, error) {
	queries, err := db.GetDBConnection(ctx, standalone, envConfig, applyMigrations, false, logger)
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
