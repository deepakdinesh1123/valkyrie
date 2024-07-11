package execution

import "github.com/google/uuid"

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Script struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LanguageDependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ExecutionEnvSpec struct {
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Packages             []string              `json:"packages"`
	EnterShell           string                `json:"enterShell"`
	Scripts              []Script              `json:"scripts"`
	Dependencies         []LanguageDependency  `json:"dependencies"`
}

type ExecutionRequest struct {
	ExecutionID uuid.UUID `json:"executionId"`
	Environment string    `json:"environment"`
	File        File      `json:"file"`
}

type ExecutionResult struct {
	ExecutionID string `json:"executionId"`
	Status      string `json:"status"`
	Result      string `json:"result"`
}
