package execution

type ExecutionRequest struct {
	File                 string               `json:"file"`
	Language             string               `json:"language,omitempty"`
	Input                string               `json:"input,omitempty"`
	LanguageDependencies []string             `json:"dependencies,omitempty"`
	LanguageVersion      string               `json:"language_version,omitempty"`
	SystemDependencies   []SystemDependency   `json:"system_dependencies,omitempty"`
	CmdLineArgs          []LanguageDependency `json:"command_line_args,omitempty"`
}

type SystemDependency struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type LanguageDependency struct {
	Name     string   `json:"name,omitempty"`
	Version  string   `json:"version,omitempty"`
	Features []string `json:"features,omitempty"` // for rust dependencies
}
