package store

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
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
	Version     string `yaml:"version"`
	Default     bool   `yaml:"default,omitempty"`
	Template    string `yaml:"template,omitempty"`
	SearchQuery string `yaml:"searchQuery"`
}

type NixPackageInfo struct {
	Attribute   string `json:"attribute"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	StorePath   string `json:"store_path"`
	Description string `json:"description"`
}

var sqliteDBPath string
var odinStoreConfig string
var genDockerConfig bool

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate odin store",
	Long:  "Generate Odin store",
	RunE:  generatePackages,
}

func init() {
	GenerateCmd.Flags().StringVarP(&sqliteDBPath, "rip-db", "r", filepath.Join(os.Getenv("XDG_DATA_HOME"), "rippkgs-index.sqlite"), "The rippkgs db to use")
	GenerateCmd.Flags().StringVarP(&odinStoreConfig, "config", "c", "", "Odin store config")
	GenerateCmd.Flags().BoolVarP(&genDockerConfig, "docker", "d", false, "Generate docker config")
}

func generatePackages(cmd *cobra.Command, args []string) error {

	sqliteClient, err := getSqliteDB()
	if err != nil {
		return logError("Could not open sqlite db", err)
	}

	defer sqliteClient.Close()

	ctx := cmd.Context()
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return logError("Error fetching environment config", err)
	}

	dbConnectionOpts := db.DBConnectionOpts(
		db.ApplyMigrations(false),
		db.IsStandalone(false),
		db.IsWorker(false),
	)

	config := logs.NewLogConfig()
	logger := logs.GetLogger(config)

	queries, err := db.GetDBConnection(ctx, envConfig, logger, dbConnectionOpts)
	if err != nil {
		return logError("Error getting DB connection", err)
	}

	err = truncateStore(ctx, queries)
	if err != nil {
		return logError("Error truncating store", err)
	}

	var storeConfig StoreConfig

	if odinStoreConfig == "" {
		err = yaml.Unmarshal(Content, &storeConfig)
		if err != nil {
			truncateStore(ctx, queries) // Truncate in case of error
			return logError("Error unmarshalling store configuration", err)
		}
	} else {
		conf, err := os.ReadFile(odinStoreConfig)
		if err != nil {
			truncateStore(ctx, queries) // Truncate in case of error
			return logError("Error reading config file", err)
		}
		err = yaml.Unmarshal(conf, &storeConfig)
		if err != nil {
			truncateStore(ctx, queries) // Truncate in case of error
			return logError("Error unmarshalling store configuration", err)
		}
	}

	err = addLanguages(ctx, queries, &storeConfig)
	if err != nil {
		truncateStore(ctx, queries) // Truncate in case of error
		return logError("Error adding languages to store", err)
	}

	err = addLanguageVersions(ctx, queries, &storeConfig)
	if err != nil {
		return logError("Error adding language versions to store", err)
	}

	err = addPackages(ctx, queries, &storeConfig, sqliteClient)
	if err != nil {
		return logError("Error adding packages to store", err)
	}

	err = queries.UpdateTextSearchVector(ctx)
	if err != nil {
		return logError("Error updating search vector", err)
	}

	return nil
}

func logError(message string, err error) error {
	if err != nil {
		logs.GetLogger(logs.NewLogConfig()).Err(err).Msg(message)
	}
	return err
}

func truncateStore(ctx context.Context, queries db.Store) error {
	if err := queries.TruncatePackages(ctx); err != nil {
		return logError("Error truncating packages", err)
	}
	if err := queries.TruncateLanguages(ctx); err != nil {
		return logError("Error truncating languages", err)
	}
	if err := queries.TruncateLanguageVersions(ctx); err != nil {
		return logError("Error truncating language versions", err)
	}
	return nil
}

func addLanguages(ctx context.Context, queries db.Store, config *StoreConfig) error {
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
		return err
	}
	return nil
}

func addLanguageVersions(ctx context.Context, queries db.Store, config *StoreConfig) error {
	var languageVersions []db.InsertLanguageVersionsParams

	for _, language := range config.Languages {
		lang, err := queries.GetLanguageByName(ctx, language.Name)
		if err != nil {
			// Log the error and continue to the next language
			logError("Error getting language by name: "+lang.Name, err)
			continue
		}
		for _, version := range language.Versions {
			langVersion := db.InsertLanguageVersionsParams{
				LanguageID:     lang.ID,
				Version:        version.Version,
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
		return logError("Error inserting language versions into the database", err)
	}
	return nil
}

func addPackages(ctx context.Context, queries db.Store, config *StoreConfig, sqliteClient *sql.DB) error {
	var packages []db.InsertPackagesParams

	for _, language := range config.Languages {
		for _, version := range language.Versions {
			pkg := db.InsertPackagesParams{
				Name:    version.NixPackage,
				Version: version.Version,
				Pkgtype: "system",
			}
			pkgInfo, err := getPackageData(pkg.Name, sqliteClient)
			if err != nil {
				return fmt.Errorf("could not get system package %s data: %s", version.NixPackage, err)
			}
			if pkgInfo.StorePath == "<broken>" {
				fmt.Printf("The language %s is broken\n", pkgInfo.Attribute)
			}
			pkg.StorePath = pgtype.Text{
				String: fmt.Sprintf("/nix/store/%s/bin", pkgInfo.StorePath),
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
					fmt.Printf("%s dependency does not exist\n", fmt.Sprintf("%s.%s", version.SearchQuery, dep))
					continue
				}
				if langPkgInfo.StorePath == "<broken>" {
					fmt.Printf("The language dependency %s is broken\n", langPkgInfo.Attribute)
					continue
				}
				langPkg.StorePath = pgtype.Text{
					String: fmt.Sprintf("/nix/store/%s", langPkgInfo.StorePath),
					Valid:  true,
				}
				langPkg.Version = langPkgInfo.Version
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
			// return fmt.Errorf("could not get system deps package data: %s", err)
			log.Println(pkg)
		}
		if sysPkgInfo.StorePath == "<broken>" {
			fmt.Printf("The system dependency %s is broken\n", sysPkgInfo.Attribute)
		}
		sysPkg.StorePath = pgtype.Text{
			String: fmt.Sprintf("/nix/store/%s/bin", sysPkgInfo.StorePath),
			Valid:  true,
		}
		sysPkg.Version = sysPkgInfo.Version
		packages = append(packages, sysPkg)
	}

	_, err := queries.InsertPackages(ctx, packages)
	if err != nil {
		return logError("Error inserting packages into the database", err)
	}

	if err := generateDockerStoreConfig(packages); err != nil {
		return logError("Error generating docker store config", err)
	}

	return nil
}

func getSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqliteDBPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getPackageData(pkg string, sqliteClient *sql.DB) (NixPackageInfo, error) {
	var nixPkgInfo NixPackageInfo
	row := sqliteClient.QueryRow("SELECT attribute, name, version, storePath, description FROM packages WHERE attribute = ?", pkg)
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
