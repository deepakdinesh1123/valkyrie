package system

import (
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

type SystemProvider struct {
	baseDir string
	cleanUp bool
	queries *db.Queries
	logger  *zerolog.Logger
}

func NewSystemProvider(baseDir string, cleanUp bool, queries *db.Queries, logger *zerolog.Logger) (*SystemProvider, error) {
	return &SystemProvider{
		baseDir: baseDir,
		cleanUp: cleanUp,
		queries: queries,
		logger:  logger,
	}, nil
}
