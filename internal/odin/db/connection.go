package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/pgembed"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed all:migrations/*.sql
var migrationsFS embed.FS

type DBOpts struct {
	standalone      bool
	applyMigrations bool
	worker          bool
	tp              trace.TracerProvider
}

type DBOptsFunc func(*DBOpts)

func IsStandalone(standalone bool) DBOptsFunc {
	return func(opts *DBOpts) {
		opts.standalone = standalone
	}
}

func ApplyMigrations(applyMigrations bool) DBOptsFunc {
	return func(opts *DBOpts) {
		opts.applyMigrations = applyMigrations
	}
}

func IsWorker(worker bool) DBOptsFunc {
	return func(opts *DBOpts) {
		opts.worker = worker
	}
}

func WithTracerProvider(tp trace.TracerProvider) DBOptsFunc {
	return func(opts *DBOpts) {
		opts.tp = tp
	}
}

func DBConnectionOpts(opts ...DBOptsFunc) *DBOpts {
	dbOpts := &DBOpts{}
	for _, opt := range opts {
		opt(dbOpts)
	}
	return dbOpts
}

func GetDBConnection(ctx context.Context, envConfig *config.EnvConfig, logger *zerolog.Logger, dbOpts *DBOpts) (Store, error) {
	// Start embedded Postgres if standalone mode is enabled
	var pge *embeddedpostgres.EmbeddedPostgres
	if dbOpts.standalone && !dbOpts.worker {
		var err error
		pge, err = pgembed.Start(
			envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_PORT,
			envConfig.POSTGRES_DB, envConfig.POSTGRES_STANDALONE_PATH)
		if err != nil {
			return nil, err
		}
	}

	// Build Postgres connection URL
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_HOST,
		envConfig.POSTGRES_PORT, envConfig.POSTGRES_DB, envConfig.POSTGRES_SSL_MODE)

	connPool, err := pgxpool.NewWithConfig(ctx, config.PgxConfig(POSTGRES_URL, dbOpts.tp, logger))
	if err != nil {
		return nil, err
	}

	// Ensure the connection is closed when the context is done
	go func() {
		<-ctx.Done()
		logger.Info().Msg("Stopping Postgres connection")
		connPool.Close()
		logger.Info().Msg("Postgres connection stopped")
		if pge != nil && !dbOpts.worker {
			logger.Info().Msg("Stopping Embedded Postgres")
			err = pge.Stop()
			if err != nil {
				logger.Err(err).Msgf("Error stopping Embedded Postgres: %s", err)
			}
			logger.Info().Msg("Embedded Postgres stopped")
		}
	}()

	// Apply migrations if requested
	if dbOpts.applyMigrations {
		logger.Info().Msg("Applying migrations")
		if err := applyMigrationsFunc(POSTGRES_URL, logger); err != nil {
			return nil, err
		}
		logger.Info().Msg("Migrations applied")
	}
	return NewStore(connPool), nil
}

// Helper function to apply migrations
func applyMigrationsFunc(postgresUrl string, logger *zerolog.Logger) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("migrations", d, postgresUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No new migrations to apply")
			return nil
		}
		return err
	}

	return nil
}
