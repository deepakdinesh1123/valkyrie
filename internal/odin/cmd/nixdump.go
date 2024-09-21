package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v38/github"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/oauth2"
)

type Package struct {
	Name    string
	Version string
}

var NixDumpCmd = &cobra.Command{
	Use:   "nixdump [channel]",
	Short: "Fetches Nixpkgs channel packages and stores them in a Postgres DB",
	RunE:  runNixDump,
}

func runNixDump(cmd *cobra.Command, args []string) error {
	channel := cmd.Flag("channel").Value.String()
	githubToken := cmd.Flag("github-token").Value.String()

	ctx := context.Background()
	pgContainer, db, err := setupTestcontainer(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup testcontainer: %w", err)
	}
	defer pgContainer.Terminate(ctx)

	client := getGithubClient(ctx, githubToken)

	branch, err := fetchReleaseByChannel(ctx, client, channel)
	if err != nil {
		return fmt.Errorf("failed to fetch branch: %w", err)
	}

	if err := processCommit(branch.GetCommit().GetSHA(), db); err != nil {
		return fmt.Errorf("failed to process commit: %w", err)
	}

	if err := createSQLDump(pgContainer, channel); err != nil {
		return fmt.Errorf("failed to create SQL dump: %w", err)
	}

	log.Println("Process completed successfully.")
	return nil
}

func init() {
	NixDumpCmd.Flags().StringP("github-token", "t", "", "GitHub token")
	NixDumpCmd.Flags().StringP("channel", "c", "", "Nixpkgs channel")
}

func setupTestcontainer(ctx context.Context) (*postgres.PostgresContainer, *sql.DB, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.5",
		postgres.WithDatabase("nixpkgs"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	connStr += " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS packages (name TEXT, version TEXT);`)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create packages table: %w", err)
	}

	return pgContainer, db, nil
}

func getGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func fetchReleaseByChannel(ctx context.Context, client *github.Client, channel string) (*github.Branch, error) {
	branchName := fmt.Sprintf("release-%s", channel)
	branch, _, err := client.Repositories.GetBranch(ctx, "NixOS", "nixpkgs", branchName, false)
	return branch, err
}

func processCommit(sha string, db *sql.DB) error {
	log.Printf("Processing commit: %s", sha)

	packages, err := extractPackages(sha)
	if err != nil {
		return fmt.Errorf("error extracting packages for commit %s: %w", sha, err)
	}

	return insertPackages(db, packages)
}

func extractPackages(sha string) ([]Package, error) {
	cmd := exec.Command("nix-env", "-qa", "--json", "-f", fmt.Sprintf("https://github.com/NixOS/nixpkgs/archive/%s.tar.gz", sha))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packagesMap map[string]map[string]interface{}
	if err := json.Unmarshal(output, &packagesMap); err != nil {
		return nil, err
	}

	packages := make([]Package, 0, len(packagesMap))
	for name, attrs := range packagesMap {
		version, _ := attrs["version"].(string)
		packages = append(packages, Package{
			Name:    strings.TrimPrefix(name, "nixpkgs."),
			Version: version,
		})
	}

	return packages, nil
}

func insertPackages(db *sql.DB, packages []Package) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO packages (name, version) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pkg := range packages {
		if _, err := stmt.Exec(pkg.Name, pkg.Version); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func createSQLDump(container *postgres.PostgresContainer, channel string) error {
	timestamp := time.Now().Format("20060102-150405") // Format: YYYYMMDD-HHMMSS
	dumpFileName := fmt.Sprintf("%s-nixos-%s.sql", timestamp, channel)
	dumpFile := filepath.Join(os.TempDir(), dumpFileName)

	ctx := context.Background()
	host, err := container.Host(ctx)
	if err != nil {
		return err
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"pg_dump",
		"-h", host,
		"-p", mappedPort.Port(),
		"-U", "user",
		"-d", "nixpkgs",
		"-f", dumpFile,
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD=password")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create SQL dump: %w", err)
	}

	fmt.Printf("SQL dump created at: %s\n", dumpFile)
	return nil
}
