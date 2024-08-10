package pgembed

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/rs/zerolog"
)

// Start starts a new Postgres database.
//
// Parameters:
// - pg_user: the Postgres username
// - pg_password: the Postgres password
// - pg_port: the Postgres port
// - pg_db: the Postgres database name
// - dataPath: the path where the database data will be stored
// - logger: the logger instance to log any errors
//
// Returns:
// - *embeddedpostgres.EmbeddedPostgres: the started Postgres database instance
// - error: any error that occurred during the database start process
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
