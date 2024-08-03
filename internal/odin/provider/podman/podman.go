package podman

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

type PodmanProvider struct {
	queries   *db.Queries
	envConfig *config.EnvConfig
	conn      context.Context
	logger    *zerolog.Logger
}

func NewPodmanProvider(env *config.EnvConfig, queries *db.Queries, logger *zerolog.Logger) (*PodmanProvider, error) {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/podman/podman.sock")
	if err != nil {
		return nil, err
	}
	return &PodmanProvider{
		queries:   queries,
		envConfig: env,
		logger:    logger,
		conn:      conn,
	}, nil
}
