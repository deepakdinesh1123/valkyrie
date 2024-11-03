package execution

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

func ConvertExecSpecToNixScript(execReq db.ExecRequest) (string, *ExecutionRequest, error) {
	execSpec := ExecutionRequest{}
	execSpec.LanguageDependencies = execReq.LanguageDependencies
	execSpec.SystemDependencies = execReq.SystemDependencies
	execSpec.CmdLineArgs = execReq.CmdLineArgs.String
	execSpec.CompileArgs = execReq.CompileArgs.String
	execSpec.Input = execReq.Input.String
	execSpec.Command = execReq.Command.String
	execSpec.Setup = execReq.Setup.String
	execSpec.Language = execReq.ProgrammingLanguage
	execSpec.LangNixPkg = config.Languages[execReq.ProgrammingLanguage]["nixPackageName"]
	scriptName := fmt.Sprintf("main.%s", config.Languages[execReq.ProgrammingLanguage]["extension"])
	execSpec.ScriptName = scriptName

	execSpec.IsFlake = false
	tmplName := filepath.Join("templates", config.Languages[execSpec.Language]["template"])

	var res bytes.Buffer
	tmpl, err := template.New(string("base.exec.tmpl")).ParseFS(
		flakes,
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
	return res.String(), &execSpec, nil
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec ExecutionRequest) (string, error) {
	execSpec.IsFlake = true
	tmplName := filepath.Join("templates", config.Languages[execSpec.Language]["template"])

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
