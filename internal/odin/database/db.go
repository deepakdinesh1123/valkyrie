package database

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func GetDBConnection(ctx context.Context, logger *zerolog.Logger) (*pgx.Conn, *db.Queries, error) {
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get env config")
		return nil, nil, err
	}
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_HOST, envConfig.POSTGRES_PORT, envConfig.POSTGRES_DB, envConfig.POSTGRES_SSL_MODE)
	DB, err := pgx.Connect(context.Background(), POSTGRES_URL)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to Postgres")
		return nil, nil, err
	}
	if envConfig.DB_MIGRATE {
		m, err := migrate.New("file://database/migrations", POSTGRES_URL)
		if err != nil {
			logger.Err(err).Msg("Failed to create migration")
			return nil, nil, err
		}
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			logger.Err(err).Msg("Failed to run migration")
			return nil, nil, err
		}
		logger.Info().Msg("Migrated database")
	}
	queries := db.New(DB)
	return DB, queries, nil
}
