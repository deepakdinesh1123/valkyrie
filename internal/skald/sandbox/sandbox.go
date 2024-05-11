package sandbox

type SandboxInterface interface {
	Create() error
	Delete() error
	SSH() error
}
