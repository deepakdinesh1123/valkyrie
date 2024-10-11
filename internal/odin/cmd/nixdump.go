package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/ulikunitz/xz"
	"golang.org/x/net/html"
)

type Package struct {
	Name     string
	Version  string
	Type     string
	Path     string
	Language string
}

var NixDumpCmd = &cobra.Command{
	Use:   "nixdump [channel]",
	Short: "Fetches Nixpkgs channel packages and stores them in a Postgres DB Dump",
	RunE:  runNixDump,
}

var (
	languagePatterns = []string{"python", "nodejs", "gcc", "rust", "php", "go", "java", "ruby"}
	systemPatterns   = []string{"glibc", "openssl", "zlib", "curl", "systemd", "sqlite"}
	packageRegex     = regexp.MustCompile(`^/nix/store/[^-]+(?:-(?P<lang>[^-]+?))?-(?P<name>[^-]+?)-(?P<version>[0-9][^/-]*)(?:-(?P<suffix>[^/-]+))?$`)
)

func runNixDump(cmd *cobra.Command, args []string) error {
	channel := cmd.Flag("channel").Value.String()

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		return fmt.Errorf("failed to get environment configuration: %w", err)
	}

	if channel == "" {
		channel = envConfig.NIXOS_VERSION
	}

	url := fmt.Sprintf("https://channels.nixos.org/nixos-%s", channel)

	doc, err := fetchHTML(url)
	if err != nil {
		return fmt.Errorf("failed to fetch HTML: %w", err)
	}

	href, err := findStorePaths(doc)
	if err != nil {
		return fmt.Errorf("failed to find the store paths href: %w", err)
	}

	if href != "" {
		if strings.HasPrefix(href, "/") {
			href = "https://releases.nixos.org" + href
		}

		finalURL, err := followRedirects(href)
		if err != nil {
			return fmt.Errorf("failed to follow redirects: %w", err)
		}

		fileName := path.Base(finalURL)
		if err := downloadFile(finalURL, fileName); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		if err := processPaths(envConfig); err != nil {
			return fmt.Errorf("failed to process paths: %w", err)
		}

		fmt.Printf("File downloaded and data inserted successfully: %s\n", fileName)

		if err := createDatabaseDump(channel, envConfig); err != nil {
			return fmt.Errorf("failed to create database dump: %w", err)
		}
	}
	return nil
}

func init() {
	NixDumpCmd.Flags().StringP("channel", "c", "", "Nixpkgs channel")
}

func processPaths(envConfig *config.EnvConfig) error {
	file, err := os.Open("store-paths.xz")
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return fmt.Errorf("could not create XZ reader: %w", err)
	}

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

	var wg sync.WaitGroup
	packageChan := make(chan Package, 100)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pkg := range packageChan {
				_, err := db.Exec(
					`INSERT INTO packages (name, version, type, language, path)
				 VALUES ($1, $2, $3, NULLIF($4, ''), $5)`,
					pkg.Name, pkg.Version, pkg.Type, pkg.Language, pkg.Path)

				if err != nil {
					log.Printf("Error inserting package: %v", err)
				}
			}
		}()
	}

	scanner := bufio.NewScanner(xzReader)
	for scanner.Scan() {
		line := scanner.Text()
		if pkg := parseLine(line); pkg != nil {
			packageChan <- *pkg
		}
	}

	close(packageChan)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading lines: %w", err)
	}

	return nil
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

func parseLine(line string) *Package {
	matches := packageRegex.FindStringSubmatch(line)
	if len(matches) > 0 {
		name := matches[packageRegex.SubexpIndex("name")]
		version := matches[packageRegex.SubexpIndex("version")]
		language := matches[packageRegex.SubexpIndex("lang")]

		if isLanguageDependency(language) {
			return &Package{
				Name:     name,
				Version:  version,
				Language: language,
				Path:     line,
				Type:     "language",
			}
		} else if isSystemDependency(name) {
			return &Package{
				Name:     name,
				Version:  version,
				Language: "",
				Path:     line,
				Type:     "system",
			}
		}
	}
	return nil
}

func isLanguageDependency(language string) bool {
	if language == "" {
		return false
	}
	for _, pattern := range languagePatterns {
		if strings.HasPrefix(language, pattern) {
			return true
		}
	}
	return false
}

func isSystemDependency(name string) bool {
	for _, pattern := range systemPatterns {
		if strings.HasPrefix(name, pattern) {
			return true
		}
	}
	return false
}

func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

func findStorePaths(doc *html.Node) (string, error) {
	var findHref func(*html.Node) string

	findHref = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "store-paths.xz") {
					return attr.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if href := findHref(c); href != "" {
				return href
			}
		}
		return ""
	}

	href := findHref(doc)

	if href == "" {
		return "", fmt.Errorf("href for store-paths.xz not found")
	}

	return href, nil
}

func followRedirects(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		return resp.Header.Get("Location"), nil
	}

	return url, nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
