package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *OdinServer) GetSystemPackages(ctx context.Context) (*api.Package, error) {
	panic("unimplemented")
}

func (s *OdinServer) SearchPackages(ctx context.Context, params api.SearchPackagesParams) (api.SearchPackagesRes, error) {
	panic("unimplemented")
}
