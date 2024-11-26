package execution

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5"
)

func ConvertExecSpecToNixScript(ctx context.Context, execReq *db.ExecRequest, queries db.Store) (string, *ExecutionRequest, error) {
	execSpec := ExecutionRequest{}
	execSpec.LanguageDependencies = execReq.LanguageDependencies
	execSpec.SystemDependencies = execReq.SystemDependencies
	execSpec.CmdLineArgs = execReq.CmdLineArgs.String
	execSpec.CompileArgs = execReq.CompileArgs.String
	execSpec.Input = execReq.Input.String
	execSpec.Command = execReq.Command.String
	execSpec.Setup = execReq.Setup.String

	langVersion, err := queries.GetLanguageVersionByID(ctx, execReq.LanguageVersion)
	if err != nil {
		return "", nil, err
	}

	language, err := queries.GetLanguageByID(ctx, langVersion.LanguageID)
	if err != nil {
		return "", nil, err
	}

	execSpec.LangNixPkg = langVersion.NixPackageName
	scriptName := fmt.Sprintf("main.%s", language.Extension)
	execSpec.ScriptName = scriptName

	execSpec.IsFlake = false

	var res bytes.Buffer
	var langTemplate string

	if langVersion.Template.String != "" {
		langTemplate = langVersion.Template.String
	} else {
		langTemplate = language.Template
	}

	langTmpl, err := template.New(string("langTmpl")).Parse(langTemplate)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse language template")
	}

	tmpl, err := langTmpl.New(string("base.exec.tmpl")).ParseFS(
		ExecTemplates,
		"templates/base.exec.tmpl",
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse template")
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute template")
	}

	fmt.Println("The nix script is\n", res.String())

	return res.String(), &execSpec, nil
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec ExecutionRequest) (string, error) {
	execSpec.IsFlake = true

	var res bytes.Buffer

	langTmpl, err := template.New(string("langTmpl")).Parse(execSpec.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse language template")
	}

	tmpl, err := langTmpl.New("base.flake.tmpl").ParseFS(ExecTemplates, "templates/base.flake.tmpl")
	if err != nil {
		s.logger.Err(err).Msg("failed to parse template")
		return "", fmt.Errorf("error parsing template: %s", err)
	}

	err = tmpl.Execute(&res, execSpec)
	if err != nil {
		s.logger.Err(err).Msg("failed to execute template")
		return "", fmt.Errorf("failed to execute template: %s", err)
	}
	return res.String(), nil
}

func (s *ExecutionService) CheckExecRequest(ctx context.Context, req *api.ExecutionRequest) error {
	language, err := s.queries.GetLanguageByName(ctx, req.Language.Value)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("specified language is not supported")
		}
	}

	if req.Version.Set {
		_, err = s.queries.GetLanguageVersion(ctx, db.GetLanguageVersionParams{
			LanguageID: language.ID,
			Version:    req.Version.Value,
		})
		if err != nil {
			return fmt.Errorf("specified version is not supported")
		}
	}

	if req.Environment.Set {
		var packages []string
		if len(req.Environment.Value.LanguageDependencies) != 0 {
			packages = append(packages, req.Environment.Value.LanguageDependencies...)
		}
		if len(req.Environment.Value.SystemDependencies) != 0 {
			packages = append(packages, req.Environment.Value.SystemDependencies...)
		}

		res, err := s.queries.PackagesExist(ctx, packages)
		if err != nil {
			return fmt.Errorf("error checking if packages exist: %s", err)
		}
		if !res.Exists {
			return fmt.Errorf("following packages does not exist: %s", res.NonexistingPackages)
		}
	}
	return nil
}
