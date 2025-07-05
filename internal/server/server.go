package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/services/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/services/sandbox"
	"github.com/deepakdinesh1123/valkyrie/internal/store"
	"github.com/deepakdinesh1123/valkyrie/internal/telemetry"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ValkyrieServer struct {
	queries          db.Store
	envConfig        *config.EnvConfig
	executionService *execution.ExecutionService
	sandboxService   *sandbox.SandboxService
	logger           *zerolog.Logger
	server           *api.Server
	tp               trace.TracerProvider
	mp               metric.MeterProvider
	prop             propagation.TextMapPropagator
	otelShutdown     func(context.Context) error
}

func NewServer(ctx context.Context, envConfig *config.EnvConfig, standalone bool, applyMigrations bool, initialiseDB bool, logger *zerolog.Logger) (*ValkyrieServer, error) {
	otelShutdown, tp, mp, prop, err := telemetry.SetupOTelSDK(ctx, "Valkyrie Server", envConfig)
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

	valkyrieServer := &ValkyrieServer{
		queries:      queries,
		envConfig:    envConfig,
		logger:       logger,
		tp:           tp,
		mp:           mp,
		otelShutdown: otelShutdown,
		prop:         prop,
	}

	if envConfig.ENABLE_EXECUTION {
		if initialiseDB {
			langs, err := queries.GetAllLanguages(ctx)
			if err != nil {
				if err == pgx.ErrNoRows {
					logger.Err(err).Msg("Generating store packages")
					store.GeneratePackages(ctx, "", "", envConfig, logger)
				}
			} else if len(langs) == 1 {
				logger.Info().Msg("Generating store packages")
				store.GeneratePackages(ctx, "", "", envConfig, logger)
			}
		}
		executionService := execution.NewExecutionService(queries, envConfig, logger)
		valkyrieServer.executionService = executionService
	}

	if envConfig.ENABLE_SANDBOX {
		sandboxService := sandbox.NewSandboxService(queries, envConfig, logger)
		valkyrieServer.sandboxService = sandboxService
	}

	srv, err := api.NewServer(
		valkyrieServer,
		api.WithTracerProvider(tp),
		api.WithMeterProvider(mp),
		api.WithPathPrefix("/api"),
	)
	if err != nil {
		return nil, err
	}

	valkyrieServer.server = srv
	return valkyrieServer, nil
}
