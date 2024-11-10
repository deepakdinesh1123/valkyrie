package execution

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
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
	tmplName := filepath.Join("templates", language.Name, langVersion.Template)

	var res bytes.Buffer
	tmpl, err := template.New(string("base.exec.tmpl")).ParseFS(
		ExecTemplates,
		"templates/base.exec.tmpl",
		tmplName,
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
	tmplName := filepath.Join("templates", execSpec.Language, execSpec.Template)

	s.logger.Debug().Str("Name", tmplName).Msg("Template is")

	var res bytes.Buffer
	tmpl, err := template.New("base.flake.tmpl").ParseFS(ExecTemplates, "templates/base.flake.tmpl", tmplName)
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
