package execution

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ExecutionRequest struct {
	ExecutionID          string   `json:"execution_id"`
	Flake                string   `json:"flake"`
	File                 File     `json:"file"`
	Language             string   `json:"language"`
	LanguageDependencies []string `json:"language_dependencies"`
	SystemDependencies   []string `json:"system_dependencies"`
	Args                 string   `json:"args"`
}
