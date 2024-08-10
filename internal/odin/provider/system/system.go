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

// NewSystemProvider creates a new SystemProvider instance.
//
// Parameters:
// - envConfig: A pointer to a config.EnvConfig object containing the environment configuration.
// - queries: A pointer to a db.Queries object containing the database queries.
// - logger: A pointer to a zerolog.Logger object for logging.
//
// Returns:
// - A pointer to a SystemProvider object if successful.
// - An error if the instance creation fails.
func NewSystemProvider(envConfig *config.EnvConfig, queries *db.Queries, logger *zerolog.Logger) (*SystemProvider, error) {
	return &SystemProvider{
		envConfig: envConfig,
		queries:   queries,
		logger:    logger,
	}, nil
}
