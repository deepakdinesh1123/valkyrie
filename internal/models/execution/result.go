package execution

type ExecutionResult struct {
	ExitCode    int    `json:"exitCode,omitempty"`
	Stderr      string `json:"stderr,omitempty"`
	Stdout      string `json:"stdout,omitempty"`
	Error       string `json:"error,omitempty"`
	ExecutionID string `json:"executionID"`
}
