package models

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

type Environment struct {
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Packages             []string              `json:"packages"`
	EnterShell           string                `json:"enterShell"`
	Scripts              []Script              `json:"scripts"`
}

type ExecutionRequest struct {
	ExecutionID string `json:"execution_id"`
	Environment string `json:"environment"`
	File        File   `json:"file"`
	Priority    int    `json:"priority"`
	Language    string `json:"language"`
	ScriptName  string `json:"script_name"`

	dockerConfig *DockerContainerConfig
	podmanConfig *PodmanContainerConfig
	systemConfig *SystemConfig
	nsjailConfig *NsjailConfig
}

type ExecutionResult struct {
	ExecutionID string `json:"execution_id"`
	Status      string `json:"status"`
	Result      string `json:"result"`
}
