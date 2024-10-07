package container

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/common"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Container struct {
	Name         string
	ID           string
	PID          int
	OverlayStore string
}

type ContainerProvider struct {
	queries   db.Store
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	user      string
}

func NewContainerExecutor(env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*ContainerProvider, error) {
	user, err := common.GetUserInfo()
	if err != nil {
		return nil, err
	}
	return &ContainerProvider{
		envConfig: env,
		logger:    logger,
		queries:   queries,
		workerId:  workerId,
		tp:        tp,
		mp:        mp,
		user:      user.Username,
	}, nil
}
