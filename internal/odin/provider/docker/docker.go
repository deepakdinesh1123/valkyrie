package docker

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
)

type DockerProvider struct {
	queries   *db.Queries
	client    *client.Client
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
}

func NewDockerProvider(env *config.EnvConfig, queries *db.Queries, logger *zerolog.Logger) (*DockerProvider, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return &DockerProvider{
		client:    client,
		envConfig: env,
		logger:    logger,
		queries:   queries,
	}, nil
}
