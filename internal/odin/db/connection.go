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

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed all:migrations/*.sql
var migrationsFS embed.FS

// GetDBConnection returns a connection to the PostgreSQL database.
//
// Parameters:
// - ctx: The context.Context used for cancellation.
// - standalone: A boolean indicating whether to start an embedded PostgreSQL instance.
// - envConfig: A pointer to the config.EnvConfig struct.
// - applyMigrations: A boolean indicating whether to apply migrations.
// - worker: A boolean indicating whether the function is being called by a worker.
// - logger: A pointer to the zerolog.Logger.
//
// Returns:
// - *Queries: A pointer to the Queries struct.
// - error: An error if any occurred.
func GetDBConnection(ctx context.Context, standalone bool, envConfig *config.EnvConfig, applyMigrations bool, worker bool, logger *zerolog.Logger) (*Queries, error) {
	// Start embedded Postgres if standalone mode is enabled
	var pge *embeddedpostgres.EmbeddedPostgres
	if standalone && !worker {
		pgDataPath := fmt.Sprintf("%s/.zango/stdb", envConfig.USER_HOME_DIR)
		var err error
		pge, err = pgembed.Start(
			envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_PORT,
			envConfig.POSTGRES_DB, pgDataPath, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to start Postgres")
			return nil, err
		}
	}

	// Build Postgres connection URL
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_HOST,
		envConfig.POSTGRES_PORT, envConfig.POSTGRES_DB, envConfig.POSTGRES_SSL_MODE)

	connPool, err := pgxpool.NewWithConfig(ctx, config.Config(POSTGRES_URL, logger))
	if err != nil {
		logger.Err(err).Msg("Failed to create connection pool")
		return nil, err
	}

	// Ensure the connection is closed when the context is done
	go func() {
		<-ctx.Done()
		logger.Info().Msg("Stopping Postgres connection")
		connPool.Close()
		logger.Info().Msg("Postgres connection stopped")
		if pge != nil && !worker {
			logger.Info().Msg("Stopping Embedded Postgres")
			err = pge.Stop()
			if err != nil {
				logger.Err(err).Msgf("Error stopping Embedded Postgres: %s", err)
			}
			logger.Info().Msg("Embedded Postgres stopped")
		}
	}()

	// Apply migrations if requested
	if applyMigrations {
		logger.Info().Msg("Applying migrations")
		if err := applyMigrationsFunc(POSTGRES_URL, logger); err != nil {
			return nil, err
		}
		logger.Info().Msg("Migrations applied")
	}

	queries := New(connPool)
	return queries, nil
}

// Helper function to apply migrations
func applyMigrationsFunc(postgresUrl string, logger *zerolog.Logger) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		logger.Err(err).Msg("Failed to create migrations")
		return err
	}
	m, err := migrate.NewWithSourceInstance("migrations", d, postgresUrl)
	if err != nil {
		logger.Err(err).Msg("Failed to create migrations instance")
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info().Msg("No new migrations to apply")
			return nil
		}
		logger.Err(err).Msg("Failed to apply migrations")
		return err
	}

	return nil
}
