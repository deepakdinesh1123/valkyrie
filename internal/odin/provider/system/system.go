package system

type SystemProvider struct {
}

func NewSystemProvider() (*SystemProvider, error) {
	return &SystemProvider{}, nil
}
