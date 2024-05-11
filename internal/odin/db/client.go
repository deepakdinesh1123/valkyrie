package db

import (
	"context"
	"fmt"

	ValkyrieConfig "github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Query *Queries

func init() {
	if Query != nil {
		return
	}
	DATABASE_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", ValkyrieConfig.EnvConfig.POSTGRES_USER, ValkyrieConfig.EnvConfig.POSTGRES_PASSWORD, ValkyrieConfig.EnvConfig.POSTGRES_HOST, ValkyrieConfig.EnvConfig.POSTGRES_PORT, ValkyrieConfig.EnvConfig.POSTGRES_DB)
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, DATABASE_URL)
	if err != nil {
		logs.Logger.Err(err)
	}
	defer conn.Close()

	queries := New(conn)
	Query = queries
}
