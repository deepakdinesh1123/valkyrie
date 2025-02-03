package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/services/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/store"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
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
	prop             propagation.TextMapPropagator
	otelShutdown     func(context.Context) error
}

func NewServer(ctx context.Context, envConfig *config.EnvConfig, standalone bool, applyMigrations bool, initialiseDB bool, logger *zerolog.Logger) (*OdinServer, error) {
	otelShutdown, tp, mp, prop, err := telemetry.SetupOTelSDK(ctx, "Odin Server", envConfig)
	if err != nil {
		logger.Err(err).Msg("Failed to setup OpenTelemetry")
		return nil, err
	}
	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(applyMigrations),
		db.IsStandalone(standalone),
		db.IsWorker(false),
		db.WithTracerProvider(tp),
	)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		return nil, err
	}

	odinServer := &OdinServer{
		queries:      queries,
		envConfig:    envConfig,
		logger:       logger,
		tp:           tp,
		mp:           mp,
		otelShutdown: otelShutdown,
		prop:         prop,
	}

	if envConfig.ODIN_ENABLE_EXECUTION {
		if initialiseDB {
			langs, err := queries.GetAllLanguages(ctx)
			if err != nil {
				if err == pgx.ErrNoRows {
					logger.Err(err).Msg("Generating store packages")
					store.GeneratePackages(ctx, "", "", false, envConfig, logger)
				}
			} else if len(langs) == 0 {
				logger.Info().Msg("Generating store packages")
				store.GeneratePackages(ctx, "", "", false, envConfig, logger)
			}
		}
		executionService := execution.NewExecutionService(queries, envConfig, logger)
		odinServer.executionService = executionService
	}
	srv, err := api.NewServer(
		odinServer,
		api.WithTracerProvider(tp),
		api.WithMeterProvider(mp),
		api.WithPathPrefix("/api"),
	)
	if err != nil {
		return nil, err
	}

	odinServer.server = srv
	return odinServer, nil
}
