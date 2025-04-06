package execution

type ExecutionRequest struct {
	Flake                string
	Code                 string
	LangNixPkg           string
	Language             string
	LangVersion          int64
	Template             string
	LanguageDependencies []string
	SystemDependencies   []string
	CmdLineArgs          string
	CompilerArgs         string
	Input                string
	Command              string
	Setup                string
	SystemSetup          string
	PkgIndex             string

	// internal variables used for converting this to flake or script
	IsFlake    bool
	ScriptName string
}
