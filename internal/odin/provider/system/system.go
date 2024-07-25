package system

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

type SystemProvider struct {
	envConfig *config.EnvConfig
	queries   *db.Queries
	logger    *zerolog.Logger
}

func NewSystemProvider(envConfig *config.EnvConfig, queries *db.Queries, logger *zerolog.Logger) (*SystemProvider, error) {
	return &SystemProvider{
		envConfig: envConfig,
		queries:   queries,
		logger:    logger,
	}, nil
}
