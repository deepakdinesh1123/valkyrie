package jsonschema

type SandboxConfig struct {
	SandboxId int64

	Memory  uint // in MegaBytes
	Storage uint // in MegaBytes
	CPU     uint
}
