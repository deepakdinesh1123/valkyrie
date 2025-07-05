package execution

type ExecutionRequest struct {
	Flake                string
	Code                 string
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
	Extension            string
	IsFlake              bool
	ScriptName           string
	NIXPKGS_REV          string
	Secrets              []byte
}
