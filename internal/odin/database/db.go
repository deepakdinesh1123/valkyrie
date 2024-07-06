package database

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func GetDBConnection(ctx context.Context, logger *zerolog.Logger) (*pgx.Conn, *db.Queries, error) {
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get env config")
		return nil, nil, err
	}
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_HOST, envConfig.POSTGRES_PORT, envConfig.POSTGRES_DB)
	DB, err := pgx.Connect(context.Background(), POSTGRES_URL)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to Postgres")
		return nil, nil, err
	}
	queries := db.New(DB)
	return DB, queries, nil
}
