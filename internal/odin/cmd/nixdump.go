package cmd

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

type Package struct {
	Name     string
	Version  string
	pkgType  string
	Language string
}

var NixDumpCmd = &cobra.Command{
	Use:   "nixdump [channel]",
	Short: "Fetches Nixpkgs channel packages and stores them in a Postgres DB Dump",
	RunE:  runNixDump,
}

var (
	languagePackageRegex = regexp.MustCompile(`^([a-zA-Z0-9]+)([0-9.]*Packages)\.([a-zA-Z0-9]+),(.+)$`)
	systemPackageRegex   = regexp.MustCompile(`^([a-zA-Z0-9-]+),(.+)$`)
	languagePatterns     = []string{"python", "nodejs", "gcc", "rust", "php", "go", "java", "ruby", "lua", "haskell"}
	systemPatterns       = []string{"glibc", "openssl", "zlib", "curl", "systemd", "sqlite", "python3", "go", "elixir"}
)

func init() {
	NixDumpCmd.Flags().StringP("channel", "c", "", "Nixpkgs channel")
}

func runNixDump(cmd *cobra.Command, args []string) error {
	channel := cmd.Flag("channel").Value.String()

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return fmt.Errorf("failed to get environment configuration: %w", err)
	}

	if channel == "" {
		channel = envConfig.NIXOS_VERSION
	}

	if err := generateNixEnvData(channel); err != nil {
		return fmt.Errorf("failed to generate Nix-env data: %w", err)
	}

	if err := processPaths(envConfig); err != nil {
		return fmt.Errorf("failed to process paths: %w", err)
	}

	if err := createDatabaseDump(channel, envConfig); err != nil {
		return fmt.Errorf("failed to create database dump: %w", err)
	}

	return nil
}

func generateNixEnvData(channel string) error {
	cmd := exec.Command("nix-env", "-qa", "--json", "-f", fmt.Sprintf("channel:nixos-%s", channel))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute nix-env command: %w", err)
	}

	// Store raw JSON data
	jsonFile, err := os.Create("nixpkgs_data.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write to JSON file: %w", err)
	}

	// Parse JSON and create CSV
	var data map[string]map[string]interface{}
	err = json.Unmarshal(output, &data)
	if err != nil {
		return fmt.Errorf("failed to parse JSON data: %w", err)
	}

	csvFile, err := os.Create("nixpkgs_data.csv")
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	for name, info := range data {
		version := info["version"].(string)
		writer.Write([]string{name, version})
	}

	return nil
}

func processPaths(envConfig *config.EnvConfig) error {
	file, err := os.Open("nixpkgs_data.csv")
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envConfig.DATABASE_HOST,
		envConfig.DATABASE_PORT,
		envConfig.POSTGRES_USER,
		envConfig.DATABASE_PASSWORD,
		envConfig.DATABASE_NAME,
		envConfig.DATABASE_SSL_MODE)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer db.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading CSV: %w", err)
		}

		if pkg := parseLine(record[0], record[1]); pkg != nil {
			_, err := db.Exec(
				`INSERT INTO packages (name, version, pkgType, language)
				VALUES ($1, $2, $3, NULLIF($4, ''))`,
				pkg.Name, pkg.Version, pkg.pkgType, pkg.Language)
			if err != nil {
				log.Printf("Error inserting package: %v\n", err)
			}
		}
	}

	return nil
}

func parseLine(name, version string) *Package {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) > 1 {
		language := parts[0]
		packageName := parts[1]
		if isLanguageDependency(language) {
			return &Package{
				Name:     packageName,
				Version:  version,
				pkgType:  "language",
				Language: language,
			}
		}
	} else if isSystemDependency(name) {
		return &Package{
			Name:     name,
			Version:  version,
			pkgType:  "system",
			Language: "",
		}
	}
	return nil
}

func isLanguageDependency(language string) bool {
	for _, pattern := range languagePatterns {
		if strings.HasPrefix(strings.ToLower(language), pattern) {
			return true
		}
	}
	return false
}

func isSystemDependency(name string) bool {
	for _, pattern := range systemPatterns {
		if strings.HasPrefix(strings.ToLower(name), pattern) {
			return true
		}
	}
	return false
}

func createDatabaseDump(channel string, envConfig *config.EnvConfig) error {
	timestamp := time.Now().Format("20060102-150405")

	dumpFileName := fmt.Sprintf("nixos-%s-%s.sql", channel, timestamp)
	dumpFilePath := fmt.Sprintf("%s/%s", envConfig.DUMP_PATH, dumpFileName)

	log.Println("Starting database dump creation...")

	dockerCmd := exec.Command("docker", "exec", envConfig.DATABASE_CONTAINER,
		"pg_dump", "-U", envConfig.POSTGRES_USER, "-d", envConfig.DATABASE_NAME, "-f", dumpFilePath)

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create database dump: %w\nOutput: %s", err, string(output))
	}
	log.Println("Database dump created inside container.")

	checkCmd := exec.Command("docker", "exec", envConfig.DATABASE_CONTAINER, "ls", "-l", dumpFilePath)
	output, err = checkCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to verify dump file in container: %w\nOutput: %s", err, string(output))
	}

	hostDumpPath := fmt.Sprintf("%s/%s", envConfig.LOCAL_DUMP_PATH, dumpFileName)
	copyCmd := exec.Command("docker", "cp",
		fmt.Sprintf("%s:%s", envConfig.DATABASE_CONTAINER, dumpFilePath),
		hostDumpPath)

	output, err = copyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy database dump to host: %w\nOutput: %s", err, string(output))
	}

	if _, err := os.Stat(hostDumpPath); os.IsNotExist(err) {
		return fmt.Errorf("dump file not found on host after copy: %w", err)
	}

	log.Printf("Database dump created and copied successfully: %s\n", hostDumpPath)

	return nil
}
