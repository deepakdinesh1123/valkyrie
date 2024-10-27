package execution

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ExecutionRequest struct {
	ExecutionID          string   `json:"execution_id"`
	Flake                string   `json:"flake"`
	NixScript            string   `json:"nix_script"`
	File                 File     `json:"file"`
	LangNixPkg           string   `json:"lang_nix_pkg"`
	Language             string   `json:"language"`
	LanguageDependencies []string `json:"language_dependencies"`
	SystemDependencies   []string `json:"system_dependencies"`
	Args                 string   `json:"args"`

	// internal variables used for converting this to flake or script
	IsFlake bool
}
