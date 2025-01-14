package store

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

//go:embed store_config.yml
var Content []byte

type StoreConfig struct {
	NixChannel string     `yaml:"nix_channel"`
	Languages  []Language `yaml:"languages"`
	Packages   []string   `yaml:"packages"`
}

type Language struct {
	Name           string            `yaml:"name"`
	Extension      string            `yaml:"extension"`
	MonacoLanguage string            `yaml:"monacoLanguage"`
	Versions       []LanguageVersion `yaml:"versions"`
	Template       string            `yaml:"template"`
	Deps           []string          `yaml:"deps"`
	DepsTemplate   string            `yaml:"depsTemplate"`
	DefaultCode    string            `yaml:"defaultCode"`
}

type LanguageVersion struct {
	NixPackage  string `yaml:"nixPackage"`
	Version     string `yaml:"version,omitempty"`
	Default     bool   `yaml:"default,omitempty"`
	Template    string `yaml:"template,omitempty"`
	SearchQuery string `yaml:"searchQuery"`
}

type NixPackageInfo struct {
	Attribute   string
	Name        sql.NullString
	Version     sql.NullString
	StorePath   sql.NullString
	Description sql.NullString
}

func GeneratePackages(ctx context.Context, odinStoreConfig string, ripPkgDBPath string, genDockerConfig bool, envConfig *config.EnvConfig, logger *zerolog.Logger) error {

	if ripPkgDBPath == "" {
		ripPkgDBPath = filepath.Join(envConfig.ODIN_INFO_DIR, "rippkgs.sqlite")

		_, err := os.Stat(ripPkgDBPath)
		if os.IsNotExist(err) {
			logger.Info().Str("path", ripPkgDBPath).Msg("Donwloading rippkgs DB-------------------")
			err = downloadRipPkgsDB(ripPkgDBPath, envConfig)
			if err != nil {
				logger.Error().Err(err).Msg("Error generating packages")
				return err
			}
		}
	}

	sqliteClient, err := getSqliteDB(ripPkgDBPath)
	if err != nil {
		logger.Error().Err(err).Msg("Could not open sqlite db")
		return err
	}

	defer sqliteClient.Close()

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

	var storeConfig StoreConfig

	if odinStoreConfig == "" {
		err = yaml.Unmarshal(Content, &storeConfig)
		if err != nil {
			truncateStore(ctx, queries, logger) // Truncate in case of error
			logger.Error().Err(err).Msg("Error unmarshalling store configuration")
			return err
		}
	} else {
		conf, err := os.ReadFile(odinStoreConfig)
		if err != nil {
			truncateStore(ctx, queries, logger) // Truncate in case of error
			logger.Error().Err(err).Msg("Error reading config file")
			return err
		}
		err = yaml.Unmarshal(conf, &storeConfig)
		if err != nil {
			truncateStore(ctx, queries, logger) // Truncate in case of error
			logger.Error().Err(err).Msg("Error unmarshalling store configuration")
			return err
		}
	}

	err = addLanguages(ctx, queries, &storeConfig, logger)
	if err != nil {
		truncateStore(ctx, queries, logger) // Truncate in case of error
		return err
	}

	err = addLanguageVersions(ctx, queries, &storeConfig, sqliteClient, logger)
	if err != nil {
		return err
	}

	err = addPackages(ctx, queries, &storeConfig, sqliteClient, genDockerConfig, logger)
	if err != nil {
		return err
	}

	err = queries.UpdateTextSearchVector(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating search vector")
		return err
	}

	return nil
}

func truncateStore(ctx context.Context, queries db.Store, logger *zerolog.Logger) error {
	if err := queries.TruncatePackages(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating packages")
		return err
	}
	if err := queries.TruncateLanguages(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating languages")
		return err
	}
	if err := queries.TruncateLanguageVersions(ctx); err != nil {
		logger.Error().Err(err).Msg("Error truncating language versions")
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

func addLanguageVersions(ctx context.Context, queries db.Store, config *StoreConfig, sqliteClient *sql.DB, logger *zerolog.Logger) error {
	var languageVersions []db.InsertLanguageVersionsParams

	for _, language := range config.Languages {
		lang, err := queries.GetLanguageByName(ctx, language.Name)
		if err != nil {
			logger.Error().Err(err).Msgf("Error getting language by name: %s", language.Name)
			continue
		}
		for _, version := range language.Versions {
			pkgData, err := getPackageData(version.NixPackage, sqliteClient)
			if err != nil {
				logger.Error().Err(err).Msgf("Could not get language %s", version.NixPackage)
			}
			langVersion := db.InsertLanguageVersionsParams{
				LanguageID:     lang.ID,
				Version:        pkgData.Version.String,
				NixPackageName: version.NixPackage,
				Template:       pgtype.Text{String: version.Template, Valid: true},
				SearchQuery:    version.SearchQuery,
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

func addPackages(ctx context.Context, queries db.Store, config *StoreConfig, sqliteClient *sql.DB, genDockerConfig bool, logger *zerolog.Logger) error {
	var packages []db.InsertPackagesParams

	for _, language := range config.Languages {
		for _, version := range language.Versions {
			pkgData, err := getPackageData(version.NixPackage, sqliteClient)
			if err != nil {
				logger.Error().Err(err).Msgf("Could not get nix data for %s", version.NixPackage)
				continue
			}
			pkg := db.InsertPackagesParams{
				Name:    version.NixPackage,
				Version: pkgData.Version.String,
				Pkgtype: "system",
			}
			pkgInfo, err := getPackageData(pkg.Name, sqliteClient)
			if err != nil {
				logger.Error().Err(err).Msgf("Could not get package %s", version.NixPackage)
			}
			if pkgInfo.StorePath.String == "<broken>" {
				logger.Warn().Msgf("The language %s is broken", pkgInfo.Attribute)
			}
			pkg.StorePath = pgtype.Text{
				String: fmt.Sprintf("/nix/store/%s/bin", pkgInfo.StorePath.String),
				Valid:  true,
			}
			packages = append(packages, pkg)
			for _, dep := range language.Deps {
				langPkg := db.InsertPackagesParams{
					Name:     dep,
					Pkgtype:  "language",
					Language: pgtype.Text{String: version.SearchQuery, Valid: true},
				}
				langPkgInfo, err := getPackageData(fmt.Sprintf("%s.%s", version.SearchQuery, dep), sqliteClient)
				if err != nil {
					logger.Warn().Msgf("%s dependency does not exist", fmt.Sprintf("%s.%s", version.SearchQuery, dep))
					continue
				}
				if langPkgInfo.StorePath.String == "<broken>" {
					logger.Warn().Msgf("The language dependency %s is broken", langPkgInfo.Attribute)
					continue
				}
				langPkg.StorePath = pgtype.Text{
					String: fmt.Sprintf("/nix/store/%s", langPkgInfo.StorePath.String),
					Valid:  true,
				}
				langPkg.Version = langPkgInfo.Version.String
				packages = append(packages, langPkg)
			}
		}
	}

	for _, pkg := range config.Packages {
		sysPkg := db.InsertPackagesParams{
			Name:    pkg,
			Pkgtype: "system",
		}
		sysPkgInfo, err := getPackageData(pkg, sqliteClient)
		if err != nil {
			logger.Warn().Msgf("Package %s not found", pkg)
			continue
		}
		if sysPkgInfo.StorePath.String == "<broken>" {
			logger.Warn().Msgf("The system dependency %s is broken", sysPkgInfo.Attribute)
		}
		sysPkg.StorePath = pgtype.Text{
			String: fmt.Sprintf("/nix/store/%s/bin", sysPkgInfo.StorePath.String),
			Valid:  true,
		}
		sysPkg.Version = sysPkgInfo.Version.String
		packages = append(packages, sysPkg)
	}

	_, err := queries.InsertPackages(ctx, packages)
	if err != nil {
		logger.Error().Err(err).Msg("Error inserting packages into the database")
		return err
	}

	if genDockerConfig {
		if err := generateDockerStoreConfig(packages); err != nil {
			logger.Error().Err(err).Msg("Error generating docker store config")
			return err
		}
	}

	return nil
}

func getSqliteDB(ripPkgDBPath string) (*sql.DB, error) {
	_, err := os.Stat(ripPkgDBPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("RipPkgsDB does not exist %s", ripPkgDBPath)
	}
	db, err := sql.Open("sqlite3", ripPkgDBPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getPackageData(pkg string, sqliteClient *sql.DB) (NixPackageInfo, error) {
	var nixPkgInfo NixPackageInfo

	query := `
		SELECT attribute, name, version, storePath, description
		FROM packages
		WHERE attribute = ?;
	`

	row := sqliteClient.QueryRow(query, pkg)
	if row == nil {
		return nixPkgInfo, fmt.Errorf("record not found: %s", pkg)
	}
	err := row.Scan(&nixPkgInfo.Attribute, &nixPkgInfo.Name, &nixPkgInfo.Version, &nixPkgInfo.StorePath, &nixPkgInfo.Description)
	if err != nil {
		return nixPkgInfo, fmt.Errorf("no rows: %s", err)
	}
	return nixPkgInfo, nil
}

func generateDockerStoreConfig(pkgs []db.InsertPackagesParams) error {
	file, err := os.Create("configs/nix/packages.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	data := [][]string{
		{"name", "type", "language"},
	}
	for _, pkg := range pkgs {
		data = append(data, []string{pkg.Name, pkg.Pkgtype, pkg.Language.String})
	}
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Write any buffered data to the underlying writer
	writer.Flush()

	// Check if any error occurred during the flush
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func downloadRipPkgsDB(ripPkgDBPath string, envConfig *config.EnvConfig) error {
	out, err := os.Create(ripPkgDBPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	resp, err := http.Get(envConfig.RIPPKGS_URL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	return nil
}
