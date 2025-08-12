package store

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

//go:embed store_config.prod.yml
var ProdConfig []byte

//go:embed store_config.dev.yml
var DevConfig []byte

type StoreConfig struct {
	Languages      []Language     `yaml:"languages"`
	SystemPackages SystemPackages `yaml:"system_packages"`
}

// Language represents an individual language entry.
type Language struct {
	Name           string    `yaml:"name"`
	Extension      string    `yaml:"extension"`
	MonacoLanguage string    `yaml:"monacoLanguage,omitempty"`
	Versions       []Version `yaml:"versions"`
	Template       string    `yaml:"template,omitempty"`
	DefaultCode    string    `yaml:"defaultCode,omitempty"`
}

// Version represents a version entry for a language.
type Version struct {
	NixPackage string `yaml:"nixPackage,omitempty"`
	Default    bool   `yaml:"default,omitempty"`
	Version    string `yaml:"version"`
	Template   string `yaml:"template,omitempty"`
}

// SystemPackages represents the system_packages section.
type SystemPackages struct {
	Exclude []string `yaml:"exclude"`
	Include []string `yaml:"include"`
}

func GeneratePackages(ctx context.Context, valkyrieStoreConfig string, ripPkgDBPath string, envConfig *config.EnvConfig, logger *zerolog.Logger) error {

	var storeConfig StoreConfig

	if valkyrieStoreConfig == "" {
		if envConfig.ENVIRONMENT == "prod" {
			err := yaml.Unmarshal(ProdConfig, &storeConfig)
			if err != nil {
				logger.Error().Err(err).Msg("Error unmarshalling store configuration")
				return err
			}
		} else if envConfig.ENVIRONMENT == "dev" {
			err := yaml.Unmarshal(DevConfig, &storeConfig)
			if err != nil {
				logger.Error().Err(err).Msg("Error unmarshalling store configuration")
				return err
			}
		} else {
			return fmt.Errorf("specified Valkyrie environment is not supported")
		}
	} else {
		conf, err := os.ReadFile(valkyrieStoreConfig)
		if err != nil {
			logger.Error().Err(err).Msg("Error reading config file")
			return err
		}
		err = yaml.Unmarshal(conf, &storeConfig)
		if err != nil {
			logger.Error().Err(err).Msg("Error unmarshalling store configuration")
			return err
		}
	}

	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(false),
		db.IsStandalone(false),
		db.IsWorker(false),
	)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting DB connection")
		return err
	}

	err = truncateStore(ctx, queries, logger)
	if err != nil {
		return err
	}

	err = addLanguages(ctx, queries, &storeConfig, logger)
	if err != nil {
		truncateStore(ctx, queries, logger) // Truncate in case of error
		return err
	}

	err = addLanguageVersions(ctx, queries, &storeConfig, logger)
	if err != nil {
		return err
	}

	err = addSysPackageFilters(ctx, queries, &storeConfig)
	if err != nil {
		return err
	}

	return nil
}

func truncateStore(ctx context.Context, queries db.Store, logger *zerolog.Logger) error {
	if err := queries.TruncateLanguages(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating languages")
		return err
	}
	if err := queries.TruncateLanguageVersions(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating language versions")
		return err
	}
	if err := queries.TruncateSysPkgFilters(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating system package filters")
		return err
	}
	return nil
}

func addLanguages(ctx context.Context, queries db.Store, config *StoreConfig, logger *zerolog.Logger) error {
	var languages []db.InsertLanguagesParams

	for _, language := range config.Languages {
		languages = append(languages, db.InsertLanguagesParams{
			Name:           language.Name,
			Extension:      language.Extension,
			MonacoLanguage: language.MonacoLanguage,
			DefaultCode:    language.DefaultCode,
			Template:       language.Template,
		})
	}

	_, err := queries.InsertLanguages(ctx, languages)
	if err != nil {
		logger.Error().Err(err).Msg("Error inserting languages")
		return err
	}
	return nil
}

func addLanguageVersions(ctx context.Context, queries db.Store, config *StoreConfig, logger *zerolog.Logger) error {
	var languageVersions []db.InsertLanguageVersionsParams

	for _, language := range config.Languages {
		lang, err := queries.GetLanguageByName(ctx, language.Name)
		if err != nil {
			logger.Error().Err(err).Msgf("Error getting language by name: %s", language.Name)
			continue
		}
		for _, version := range language.Versions {
			langVersion := db.InsertLanguageVersionsParams{
				LanguageID:     lang.ID,
				Version:        version.Version,
				NixPackageName: pgtype.Text{String: version.NixPackage, Valid: true},
				Template:       pgtype.Text{String: version.Template, Valid: true},
			}
			if version.Default {
				langVersion.DefaultVersion = true
			}
			languageVersions = append(languageVersions, langVersion)
		}
	}
	_, err := queries.InsertLanguageVersions(ctx, languageVersions)
	if err != nil {
		logger.Error().Err(err).Msg("Error inserting language versions into the database")
		return err
	}
	return nil
}

func addSysPackageFilters(ctx context.Context, queries db.Store, config *StoreConfig) error {
	include := config.SystemPackages.Include
	exclude := config.SystemPackages.Exclude

	filters := []db.BulkInsertSystemPackageFiltersParams{}
	for _, regxp := range include {
		filters = append(filters, db.BulkInsertSystemPackageFiltersParams{
			FilterType:    "include",
			PackageString: regxp,
		})
	}
	for _, regxp := range exclude {
		filters = append(filters, db.BulkInsertSystemPackageFiltersParams{
			FilterType:    "exclude",
			PackageString: regxp,
		})
	}
	_, err := queries.BulkInsertSystemPackageFilters(ctx, filters)
	return fmt.Errorf("error inserting system package filters: %v", err)
}
