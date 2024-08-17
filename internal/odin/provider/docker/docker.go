package docker

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DockerProvider struct {
	queries   *db.Queries
	store     db.Store
	connPool  *pgxpool.Pool
	client    *client.Client
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
}

func NewDockerProvider(env *config.EnvConfig, queries *db.Queries, tp trace.TracerProvider, mp metric.MeterProvider, connPool *pgxpool.Pool, logger *zerolog.Logger) (*DockerProvider, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return &DockerProvider{
		client:    client,
		connPool:  connPool,
		envConfig: env,
		logger:    logger,
		queries:   queries,
		store:     db.NewStore(connPool),
		tp:        tp,
		mp:        mp,
	}, nil
}
