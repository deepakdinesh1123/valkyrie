package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
)

func (s *ValkyrieServer) GetVersion(ctx context.Context) (api.GetVersionRes, error) {
	return &api.GetVersionOK{
		Version: config.VERSION,
	}, nil
}
