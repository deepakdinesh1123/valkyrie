package execution

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
)

type ExecutionRequest struct {
	Flake                string
	Code                 string
	LangNixPkg           string
	Language             string
	LanguageDependencies []string
	SystemDependencies   []string
	CmdLineArgs          string
	CompileArgs          string
	Input                string
	Command              string
	Setup                string

	// internal variables used for converting this to flake or script
	IsFlake    bool
	ScriptName string
}

//go:embed templates
var flakes embed.FS

type ExecutionService struct {
	queries   db.Store
	envConfig *config.EnvConfig
	logger    *zerolog.Logger
}

func NewExecutionService(queries db.Store, envConfig *config.EnvConfig, logger *zerolog.Logger) *ExecutionService {
	return &ExecutionService{
		queries:   queries,
		envConfig: envConfig,
		logger:    logger,
	}
}

func (s *ExecutionService) prepareExecutionRequest(req *api.ExecutionRequest) (*ExecutionRequest, error) {
	execReq := ExecutionRequest{}
	if req.Timeout.Value > 180 {
		return nil, &ExecutionServiceError{
			Type:    "timeout",
			Message: "Timeout cannot be more than 180 seconds",
		}
	}

	languageData, err := s.queries.GetAllLanguages(context.TODO())
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "db",
			Message: "Failed to fetch languages from the database",
		}
	}

	languages := make(map[string]api.LanguageResponse)
	for _, lang := range languageData {
		languages[lang.Name] = api.LanguageResponse{
			ID:          lang.ID,
			Name:        lang.Name,
			Extension:   lang.Extension,
			DefaultCode: lang.DefaultCode,
		}
	}

	if _, ok := languages[req.Language]; !ok {
		return nil, &ExecutionServiceError{
			Type:    "language",
			Message: "Language not supported",
		}
	}

	languageVersionData, err := s.queries.GetVersionsByLanguageID(context.TODO(), languages[req.Language].ID)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "db",
			Message: "Failed to fetch language version from the database",
		}
	}

	languageVersions := make(map[string]api.LanguageVersionResponse)
	for _, lang := range languageVersionData {
		languageVersions[lang.Version] = api.LanguageVersionResponse{
			ID:             lang.ID,
			LanguageID:     lang.LanguageID,
			Version:        lang.Version,
			NixPackageName: lang.NixPackageName,
			FlakeTemplate:  lang.FlakeTemplate,
			ScriptTemplate: lang.ScriptTemplate,
			SearchQuery:    lang.SearchQuery,
		}
	}

	if _, ok := languageVersions[req.Version.Value]; !ok {
		return nil, &ExecutionServiceError{
			Type:    "language",
			Message: "Language version not supported",
		}
	}

	languageConfig := languages[req.Language]
	languageVersionConfig := languageVersions[req.Version.Value]
	execReq.Language = req.Language
	execReq.LangNixPkg = languageVersions[req.Version.Value].NixPackageName
	scriptName := fmt.Sprintf("main.%s", languageConfig.Extension)
	execReq.File = File{
		Name:    scriptName,
		Content: req.Code,
	}
	execReq.Args = req.Environment.Value.Args.Value
	execReq.SystemDependencies = req.Environment.Value.SystemDependencies
	execReq.LanguageDependencies = req.Environment.Value.LanguageDependencies

	flake, err := s.convertExecSpecToFlake(execReq, languageVersionConfig)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "flake",
			Message: err.Error(),
		}
	}

	nixScript, err := s.convertExecSpecToNixScript(execReq, languageVersionConfig)
	if err != nil {
		return nil, &ExecutionServiceError{
			Type:    "nix script",
			Message: err.Error(),
		}
	}
	execReq.Flake = flake
	execReq.NixScript = nixScript
	return &execReq, nil
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec ExecutionRequest, language api.LanguageVersionResponse) (string, error) {
	execSpec.IsFlake = true
	tmplName := filepath.Join("templates", language.FlakeTemplate)

	var res bytes.Buffer
	tmpl, err := template.New("base.flake.tmpl").ParseFS(flakes, "templates/base.flake.tmpl", tmplName)
	if err != nil {
		s.logger.Err(err).Msg("failed to parse template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		s.logger.Err(err).Msg("failed to execute template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return res.String(), nil
}

func (s *ExecutionService) convertExecSpecToNixScript(execSpec ExecutionRequest, language api.LanguageVersionResponse) (string, error) {
	execSpec.IsFlake = false
	tmplName := filepath.Join("templates", language.ScriptTemplate)

	var res bytes.Buffer
	tmpl, err := template.New(string("base.exec.tmpl")).ParseFS(
		flakes,
		"templates/base.exec.tmpl",
		tmplName,
	)
	if err != nil {
		s.logger.Err(err).Msg("failed to parse template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to parse template",
		}
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		s.logger.Err(err).Msg("failed to execute template")
		return "", &ExecutionServiceError{
			Type:    "template",
			Message: "failed to execute template",
		}
	}
	return res.String(), nil
}

func (s *ExecutionService) AddJob(ctx context.Context, req *api.ExecutionRequest) (int64, error) {
	execReq, err := s.prepareExecutionRequest(req)
	if err != nil {
		return 0, err
	}
	var jobParams db.AddJobTxParams
	jobParams.Code = execReq.File.Content
	jobParams.Flake = execReq.Flake
	jobParams.NixScript = execReq.NixScript
	jobParams.ProgrammingLanguage = req.Language
	jobParams.MaxRetries = req.MaxRetries.Value
	jobParams.Path = execReq.File.Name
	jobParams.Timeout = req.Timeout.Value

	hash := calculateHash(jobParams.Code, jobParams.ProgrammingLanguage, jobParams.Flake, jobParams.NixScript, jobParams.Path)
	jobParams.Hash = hash

	job, err := s.queries.AddJobTx(ctx, jobParams)
	if err != nil {
		s.logger.Err(err).Msg("failed to add job")
		return 0, err
	}
	return int64(job.JobID), nil
}

func calculateHash(code string, language string, flake string, nix_script string, path string) string {
	hashInstance := sha256.New()
	hashInstance.Write([]byte(code))
	hashInstance.Write([]byte(language))
	hashInstance.Write([]byte(flake))
	hashInstance.Write([]byte(path))
	hashInstance.Write([]byte(nix_script))
	return hex.EncodeToString(hashInstance.Sum(nil))
}
