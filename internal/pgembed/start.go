package pgembed

import (
	"os"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/rs/zerolog"
)

func Start(config *config.EnvConfig, logger *zerolog.Logger) (*embeddedpostgres.EmbeddedPostgres, error) {

	homedir, err := os.UserHomeDir()
	if err != nil {
		logger.Err(err).Msg("Failed to get user home directory")
		return nil, err
	}
	dataPath := homedir + "/.valkyrie/postgres/data"
	err = os.MkdirAll(dataPath, 0755)
	if err != nil {
		logger.Err(err).Msg("Failed to create data path")
		return nil, err
	}

	pg := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username(config.POSTGRES_USER).
			Password(config.POSTGRES_PASSWORD).
			Port(config.POSTGRES_PORT).
			Database(config.POSTGRES_DB).
			DataPath(dataPath),
	)
	err = pg.Start()
	if err != nil {
		logger.Err(err).Msg("Failed to start Postgres")
		return nil, err
	}
	return pg, nil
}
