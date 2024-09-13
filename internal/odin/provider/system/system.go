package system

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type SystemProvider struct {
	envConfig *config.EnvConfig
	queries   db.Store
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	workerId  int32
}

func NewSystemProvider(envConfig *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*SystemProvider, error) {
	return &SystemProvider{
		envConfig: envConfig,
		queries:   queries,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}
