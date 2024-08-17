package system

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type SystemProvider struct {
	envConfig *config.EnvConfig
	queries   *db.Queries
	connPool  *pgxpool.Pool
	logger    *zerolog.Logger
	store     db.Store
	tp        trace.TracerProvider
	mp        metric.MeterProvider
}

func NewSystemProvider(envConfig *config.EnvConfig, queries *db.Queries, tp trace.TracerProvider, mp metric.MeterProvider, connPool *pgxpool.Pool, logger *zerolog.Logger) (*SystemProvider, error) {
	return &SystemProvider{
		envConfig: envConfig,
		queries:   queries,
		connPool:  connPool,
		logger:    logger,
		store:     db.NewStore(connPool),
		tp:        tp,
		mp:        mp,
	}, nil
}
