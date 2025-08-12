package pool

type Pod struct {
	Name      string
	Container Container
	SandboxID string
}
