package pgembed

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/rs/zerolog"
)

func Start(pg_user string, pg_password string, pg_port uint32, pg_db string, dataPath string, logger *zerolog.Logger) (*embeddedpostgres.EmbeddedPostgres, error) {
	pg := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Username(pg_user).
			Password(pg_password).
			Port(pg_port).
			Database(pg_db).
			DataPath(dataPath),
	)
	err := pg.Start()
	if err != nil {
		logger.Err(err).Msg("Failed to start Postgres")
		return nil, err
	}
	return pg, nil
}
