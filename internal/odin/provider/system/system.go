package system

type SystemProvider struct {
	baseDir string
	cleanUp bool
}

func NewSystemProvider(baseDir string, cleanUp bool) (*SystemProvider, error) {
	return &SystemProvider{
		baseDir: baseDir,
		cleanUp: cleanUp,
	}, nil
}
