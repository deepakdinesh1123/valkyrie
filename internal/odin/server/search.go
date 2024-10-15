package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

// SearchLanguagePackages implements api.Handler.
func (s *OdinServer) SearchLanguagePackages(ctx context.Context, params api.SearchLanguagePackagesParams) (api.SearchLanguagePackagesRes, error) {
	dbParams := db.SearchLanguagePackagesParams{
		Language:    params.Language,
		Searchquery: params.SearchString,
	}

	packages, err := s.queries.SearchLanguagePackages(ctx, dbParams)
	if err != nil {
		return &api.SearchLanguagePackagesBadRequest{}, err
	}

	results := make([]api.Package, len(packages))
	for i, pkg := range packages {
		results[i] = api.Package{
			Name:    pkg.Name,
			Version: pkg.Version,
		}
	}

	return &api.SearchLanguagePackagesOK{
		Packages: results,
	}, nil
}

// SearchSystemPackages implements api.Handler.
func (s *OdinServer) SearchSystemPackages(ctx context.Context, params api.SearchSystemPackagesParams) (api.SearchSystemPackagesRes, error) {

	packages, err := s.queries.SearchSystemPackages(ctx, params.SearchString)
	if err != nil {
		return &api.SearchSystemPackagesBadRequest{}, err
	}

	results := make([]api.Package, len(packages))
	for i, pkg := range packages {
		results[i] = api.Package{
			Name:    pkg.Name,
			Version: pkg.Version,
		}
	}

	return &api.SearchSystemPackagesOK{
		Packages: results,
	}, nil
}
