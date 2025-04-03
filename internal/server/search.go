package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
)

// SearchLanguagePackages implements api.Handler.
func (s *ValkyrieServer) SearchLanguagePackages(ctx context.Context, params api.SearchLanguagePackagesParams) (api.SearchLanguagePackagesRes, error) {
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
func (s *ValkyrieServer) SearchSystemPackages(ctx context.Context, params api.SearchSystemPackagesParams) (api.SearchSystemPackagesRes, error) {

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

// PackagesExist implements api.Handler.
func (s *ValkyrieServer) PackagesExist(ctx context.Context, req *api.PackageExistRequest, params api.PackagesExistParams) (api.PackagesExistRes, error) {
	dbParams := db.PackagesExistParams{
		Language: req.Language,
		Packages: req.Packages,
	}

	existsResponse, err := s.queries.PackagesExist(ctx, dbParams)
	if err != nil {
		return &api.PackagesExistBadRequest{}, err
	}

	response := api.PackagesExistOK{
		Exists:              existsResponse.Exists,
		NonExistingPackages: existsResponse.NonexistingPackages,
	}

	return &response, nil
}

// FetchLanguagePackages implements api.Handler.
func (s *ValkyrieServer) FetchLanguagePackages(ctx context.Context, params api.FetchLanguagePackagesParams) (api.FetchLanguagePackagesRes, error) {
	packages, err := s.queries.FetchLanguagePackages(ctx, params.Language)
	if err != nil {
		return &api.FetchLanguagePackagesBadRequest{}, err
	}

	results := make([]api.Package, len(packages))
	for i, pkg := range packages {
		results[i] = api.Package{
			Name:    pkg.Name,
			Version: pkg.Version,
		}
	}

	return &api.FetchLanguagePackagesOK{
		Packages: results,
	}, nil
}

// FetchSystemPackages implements api.Handler.
func (s *ValkyrieServer) FetchSystemPackages(ctx context.Context, params api.FetchSystemPackagesParams) (api.FetchSystemPackagesRes, error) {
	packages, err := s.queries.FetchSystemPackages(ctx)
	if err != nil {
		return &api.FetchSystemPackagesBadRequest{}, err
	}

	results := make([]api.Package, len(packages))
	for i, pkg := range packages {
		results[i] = api.Package{
			Name:    pkg.Name,
			Version: pkg.Version,
		}
	}
	return &api.FetchSystemPackagesOK{
		Packages: results,
	}, nil
}
