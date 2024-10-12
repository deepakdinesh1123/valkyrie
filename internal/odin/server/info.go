package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *OdinServer) GetVersion(ctx context.Context, params api.GetVersionParams) (api.GetVersionRes, error) {
	return &api.GetVersionOK{
		Version: config.ODIN_VERSION,
	}, nil
}
