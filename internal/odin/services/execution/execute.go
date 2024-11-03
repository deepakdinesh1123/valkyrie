package execution

type ExecutionRequest struct {
	Flake                string
	Code                 string
	LangNixPkg           string
	Language             string
	LanguageDependencies []string
	SystemDependencies   []string
	CmdLineArgs          string
	CompileArgs          string
	Input                string
	Command              string
	Setup                string

	// internal variables used for converting this to flake or script
	IsFlake    bool
	ScriptName string
}
