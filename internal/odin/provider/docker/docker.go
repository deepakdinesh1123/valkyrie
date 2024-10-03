package docker

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/common"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DockerProvider struct {
	queries   db.Store
	client    *client.Client
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	user      string
}

func NewDockerProvider(env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*DockerProvider, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	user, err := common.GetUserInfo()
	if err != nil {
		return nil, err
	}
	return &DockerProvider{
		client:    client,
		envConfig: env,
		logger:    logger,
		queries:   queries,
		workerId:  workerId,
		tp:        tp,
		mp:        mp,
		user:      user.Username,
	}, nil
}
