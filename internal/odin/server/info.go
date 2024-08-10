package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

// GetVersion returns the version of the Odin server.
//
// Parameters:
// - ctx: the context of the request.
// Returns:
// - api.GetVersionRes: the version of the Odin server.
// - error: any error that occurred while retrieving the version.
func (s *OdinServer) GetVersion(ctx context.Context) (api.GetVersionRes, error) {
	return &api.GetVersionOK{
		Version: config.ODIN_VERSION,
	}, nil
}
