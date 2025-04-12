package execution

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/jackc/pgx/v5"
)

func ConvertExecSpecToNixScript(ctx context.Context, execReq *db.ExecRequest, queries db.Store) (string, *ExecutionRequest, error) {
	execSpec := ExecutionRequest{}
	execSpec.LanguageDependencies = execReq.LanguageDependencies
	execSpec.SystemDependencies = execReq.SystemDependencies
	execSpec.CmdLineArgs = execReq.CmdLineArgs.String
	execSpec.CompilerArgs = execReq.CompileArgs.String
	execSpec.Input = execReq.Input.String
	execSpec.Command = execReq.Command.String
	execSpec.Setup = execReq.Setup.String
	execSpec.SystemSetup = execReq.SystemSetup.String
	execSpec.PkgIndex = execReq.PkgIndex.String

	fmt.Println(execReq.LanguageDependencies, execReq.SystemDependencies, execReq.Setup.String, execReq.Command.String)

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

	baseTemplate, err := ExecTemplates.ReadFile("templates/base.exec.tmpl")
	if err != nil {
		return "", nil, fmt.Errorf("failed to read base template")
	}

	baseTmpl, err := template.New("base").Parse(string(baseTemplate))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse template")
	}

	langTmpl, err := template.Must(baseTmpl.Clone()).Parse(langTemplate)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse language template")
	}

	err = langTmpl.Execute(&res, execSpec)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute template")
	}

	fmt.Println("The nix script is\n", res.String())

	return res.String(), &execSpec, nil
}

func (s *ExecutionService) convertExecSpecToFlake(execSpec ExecutionRequest) (string, error) {
	execSpec.IsFlake = true

	var res bytes.Buffer

	baseTemplate, err := ExecTemplates.ReadFile("templates/base.flake.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read base template")
	}

	baseTmpl, err := template.New("base").Parse(string(baseTemplate))
	if err != nil {
		return "", fmt.Errorf("failed to parse template")
	}

	langTmpl, err := template.Must(baseTmpl.Clone()).Parse(execSpec.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse language template")
	}

	err = langTmpl.Execute(&res, execSpec)
	if err != nil {
		return "", fmt.Errorf("failed to execute template")
	}
	return res.String(), nil
}

func (s *ExecutionService) CheckExecRequest(ctx context.Context, req *api.ExecutionRequest) ([]string, error) {
	language, err := s.queries.GetLanguageByName(ctx, req.Language.Value)
	if err != nil {
		if err == pgx.ErrNoRows {
			langs, err := s.queries.GetAllLanguages(ctx)
			if err != nil {
				return []string{}, fmt.Errorf("error fetching all the supported languages %v", err)
			}
			langNames := []string{}
			for _, lang := range langs {
				langNames = append(langNames, lang.Name)
			}
			return langNames, fmt.Errorf("specified language is not supported")
		}
	}

	if req.Version.Set && !(req.Version.Value == "") {
		_, err = s.queries.GetLanguageVersion(ctx, db.GetLanguageVersionParams{
			LanguageID: language.ID,
			Version:    req.Version.Value,
		})
		if err != nil {
			versions, err := s.queries.GetVersionsByLanguageID(ctx, language.ID)
			if err != nil {
				return []string{}, fmt.Errorf("error fetching all the supported versions %v", err)
			}
			vers := []string{}
			for _, ver := range versions {
				vers = append(vers, ver.Version)
			}

			return vers, fmt.Errorf("specified version is not supported")
		}
	}

	// if req.Environment.Set {
	// 	var packages []string
	// 	if len(req.Environment.Value.LanguageDependencies) != 0 {
	// 		packages = append(packages, req.Environment.Value.LanguageDependencies...)
	// 	}
	// 	if len(req.Environment.Value.SystemDependencies) != 0 {
	// 		packages = append(packages, req.Environment.Value.SystemDependencies...)
	// 	}

	// 	if len(packages) > 0 {
	// 		res, err := s.queries.PackagesExist(ctx, db.PackagesExistParams{
	// 			Packages: packages,
	// 			Language: req.Language.Value,
	// 		})
	// 		if err != nil {
	// 			return fmt.Errorf("error checking if packages exist: %s", err)
	// 		}
	// 		if !res.Exists {
	// 			return fmt.Errorf("following packages does not exist: %s", res.NonexistingPackages)
	// 		}
	// 	}
	// }
	return []string{}, nil
}
