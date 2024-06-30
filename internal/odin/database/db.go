package database

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn
var Queries *db.Queries

func init() {
	if DB != nil && Queries != nil {
		return
	}
	var err error
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.EnvConfig.POSTGRES_USER, config.EnvConfig.POSTGRES_PASSWORD, config.EnvConfig.POSTGRES_HOST, config.EnvConfig.POSTGRES_PORT, config.EnvConfig.POSTGRES_DB)
	logs.Logger.Info().Msg(POSTGRES_URL)
	DB, err = pgx.Connect(context.Background(), POSTGRES_URL)
	if err != nil {
		logs.Logger.Err(err).Msg("Failed to connect to Postgres")
		panic(err)
	}
	Queries = db.New(DB)
}
