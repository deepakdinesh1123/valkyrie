package sandbox

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db/jsonschema"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
)

//go:embed templates
var SandboxTemplates embed.FS

type SandboxService struct {
	queries   db.Store
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
}

func NewSandboxService(queries db.Store, envConfig *config.EnvConfig, logger *zerolog.Logger) *SandboxService {
	return &SandboxService{
		queries:   queries,
		envConfig: envConfig,
		logger:    logger,
	}
}

func (s *SandboxService) AddSandbox(ctx context.Context, sandboxReq *api.OptCreateSandbox) (int64, error) {
	flake, err := s.generateFlake(sandboxReq)
	if err != nil {
		return 0, fmt.Errorf("failed to generate flake: %w", err)
	}

	res, err := s.queries.AddSandboxJobTx(ctx, db.AddSandboxTxParams{
		SandboxConfig: jsonschema.SandboxConfig{
			Flake: flake,
		},
		MaxRetries: s.envConfig.ODIN_MAX_RETRIES,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to add sandbox job: %w", err)
	}

	return res.SandboxId, nil
}

func (s *SandboxService) generateFlake(sandboxReq *api.OptCreateSandbox) (string, error) {
	// If sandbox request is not set, return empty flake
	if !sandboxReq.Set {
		return "", nil
	}

	// If NixFlake value is provided, use it directly
	if sandboxReq.Value.NixFlake.Value != "" {
		return sandboxReq.Value.NixFlake.Value, nil
	}

	// Generate flake from template
	flakeTemplateFile, err := SandboxTemplates.ReadFile("templates/flake.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read flake template: %w", err)
	}

	tmpl, err := template.New("flake.nix").Parse(string(flakeTemplateFile))
	if err != nil {
		return "", fmt.Errorf("failed to parse flake template: %w", err)
	}

	flakeConfig := FlakeConfig{
		NIXPKGS_URL:        fmt.Sprintf("path:%s", s.envConfig.ODIN_SANDBOX_NIXPKGS_PATH),
		SystemDependencies: sandboxReq.Value.SystemDependencies,
		Languages:          sandboxReq.Value.Languages,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, flakeConfig); err != nil {
		return "", fmt.Errorf("failed to execute flake template: %w", err)
	}

	return buf.String(), nil
}
