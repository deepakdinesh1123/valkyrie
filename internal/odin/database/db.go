package database

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/pgembed"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed all:migrations/*.sql
var migrationsFS embed.FS

func GetDBConnection(ctx context.Context, standalone bool, envConfig *config.EnvConfig, applyMigrations bool, sigChan chan os.Signal, done chan bool, logger *zerolog.Logger) (*pgx.Conn, *db.Queries, error) {
	// Start embedded Postgres if standalone mode is enabled
	var pge *embeddedpostgres.EmbeddedPostgres
	if standalone {
		var err error
		pge, err = pgembed.Start(envConfig, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to start Postgres")
			return nil, nil, err
		}
	}

	// Build Postgres connection URL
	POSTGRES_URL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		envConfig.POSTGRES_USER, envConfig.POSTGRES_PASSWORD, envConfig.POSTGRES_HOST,
		envConfig.POSTGRES_PORT, envConfig.POSTGRES_DB, envConfig.POSTGRES_SSL_MODE)

	// Connect to Postgres
	DB, err := pgx.Connect(ctx, POSTGRES_URL)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to Postgres")
		return nil, nil, err
	}

	// Ensure the connection is closed when the context is done
	go func() {
		<-sigChan
		if pge != nil {
			logger.Info().Msg("Stopping Embedded Postgres")
			err = pge.Stop()
			if err != nil {
				logger.Err(err).Msgf("Error stopping Embedded Postgres: %s", err)
			}
			logger.Info().Msg("Embedded Postgres stopped")
		}
		logger.Info().Msg("Stopping Postgres connection")
		DB.Close(ctx)
		logger.Info().Msg("Postgres connection stopped")
		done <- true
	}()

	// Apply migrations if requested
	if applyMigrations {
		logger.Info().Msg("Applying migrations")
		if err := applyMigrationsFunc(POSTGRES_URL, logger); err != nil {
			return nil, nil, err
		}
	}

	queries := db.New(DB)
	return DB, queries, nil
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
