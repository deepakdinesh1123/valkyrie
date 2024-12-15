package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateLanguage implements api.Handler.
func (s *OdinServer) CreateLanguage(ctx context.Context, req *api.Language, params api.CreateLanguageParams) (api.CreateLanguageRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.CreateLanguageForbidden{}, nil
	}

	if req == nil {
		return &api.CreateLanguageBadRequest{}, errors.New("invalid request: language is required")
	}

	dbParams := db.CreateLanguageParams{
		Name:           req.Name,
		Extension:      req.Extension,
		MonacoLanguage: req.MonacoLanguage,
		DefaultCode:    req.DefaultCode,
	}

	id, err := s.queries.CreateLanguage(ctx, dbParams)
	if err != nil {
		return &api.CreateLanguageBadRequest{}, err
	}

	return &api.CreateLanguageOK{
		Language: api.LanguageResponse{
			ID:             id,
			Name:           req.Name,
			Extension:      req.Extension,
			MonacoLanguage: req.MonacoLanguage,
			DefaultCode:    req.DefaultCode,
		},
	}, nil
}

// CreateLanguageVersion implements api.Handler.
func (s *OdinServer) CreateLanguageVersion(ctx context.Context, req *api.LanguageVersion, params api.CreateLanguageVersionParams) (api.CreateLanguageVersionRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.CreateLanguageVersionForbidden{}, nil
	}

	if req == nil {
		return &api.CreateLanguageVersionBadRequest{}, errors.New("invalid request: language is required")
	}

	dbParams := db.CreateLanguageVersionParams{
		LanguageID:     req.LanguageID,
		Version:        req.Version,
		NixPackageName: req.NixPackageName,
		Template:       pgtype.Text{String: req.Template, Valid: true},
		SearchQuery:    req.SearchQuery,
		DefaultVersion: req.DefaultVersion,
	}

	id, err := s.queries.CreateLanguageVersion(ctx, dbParams)
	if err != nil {
		return &api.CreateLanguageVersionBadRequest{}, err
	}

	return &api.CreateLanguageVersionOK{
		Language: api.LanguageVersionResponse{
			ID:             id,
			LanguageID:     req.LanguageID,
			Version:        req.Version,
			NixPackageName: req.NixPackageName,
			Template:       req.Template,
			SearchQuery:    req.SearchQuery,
			DefaultVersion: req.DefaultVersion,
		},
	}, nil
}

// DeleteLanguage implements api.Handler.
func (s *OdinServer) DeleteLanguage(ctx context.Context, params api.DeleteLanguageParams) (api.DeleteLanguageRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.DeleteLanguageForbidden{}, nil
	}

	id, err := s.queries.DeleteLanguage(ctx, params.ID)
	if id == 0 {
		return &api.DeleteLanguageBadRequest{
			Message: fmt.Sprintf("Language with id %d not found: %v", params.ID, err),
		}, nil
	}
	if err != nil {
		return &api.DeleteLanguageInternalServerError{
			Message: fmt.Sprintf("Failed to delete language: %v", err),
		}, err
	}
	return &api.DeleteLanguageOK{
		Message: fmt.Sprintf("Language with id %d deleted", params.ID),
	}, nil
}

// DeleteLanguageVersion implements api.Handler.
func (s *OdinServer) DeleteLanguageVersion(ctx context.Context, params api.DeleteLanguageVersionParams) (api.DeleteLanguageVersionRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.DeleteLanguageVersionForbidden{}, nil
	}

	id, err := s.queries.DeleteLanguageVersion(ctx, params.ID)
	if id == 0 {
		return &api.DeleteLanguageVersionBadRequest{
			Message: fmt.Sprintf("Language version with id %d not found: %v", params.ID, err),
		}, nil
	}
	if err != nil {
		return &api.DeleteLanguageVersionInternalServerError{
			Message: fmt.Sprintf("Failed to delete language version: %v", err),
		}, err
	}
	return &api.DeleteLanguageVersionOK{
		Message: fmt.Sprintf("Language version with id %d deleted", params.ID),
	}, nil
}

// GetAllLanguages implements api.Handler.
func (s *OdinServer) GetAllLanguages(ctx context.Context, params api.GetAllLanguagesParams) (api.GetAllLanguagesRes, error) {
	languages, err := s.queries.GetAllLanguages(ctx)
	if err != nil {
		s.logger.Printf("Error fetching languages: %v", err)
		return &api.GetAllLanguagesInternalServerError{
			Message: "An internal error occurred while fetching languages. Please try again later.",
		}, err
	}

	if len(languages) == 0 {
		return &api.GetAllLanguagesOK{
			Languages: []api.LanguageResponse{},
		}, nil
	}

	var languageResponses []api.LanguageResponse
	for _, lang := range languages {
		languageResponses = append(languageResponses, api.LanguageResponse{
			ID:             lang.ID,
			Name:           lang.Name,
			Extension:      lang.Extension,
			MonacoLanguage: lang.MonacoLanguage,
			DefaultCode:    lang.DefaultCode,
		})
	}

	return &api.GetAllLanguagesOK{
		Languages: languageResponses,
	}, nil
}

// GetAllLanguageVersions implements api.Handler.
func (s *OdinServer) GetAllLanguageVersions(ctx context.Context, params api.GetAllLanguageVersionsParams) (api.GetAllLanguageVersionsRes, error) {
	languageVersions, err := s.queries.GetAllLanguageVersions(ctx)
	if err != nil {
		s.logger.Printf("Error fetching language versions: %v", err)

		return &api.GetAllLanguageVersionsInternalServerError{
			Message: "An internal error occurred while fetching language versions. Please try again later.",
		}, err
	}

	if len(languageVersions) == 0 {
		return &api.GetAllLanguageVersionsOK{
			LanguageVersions: []api.LanguageVersionResponse{},
		}, nil
	}

	var languageVersionResponses []api.LanguageVersionResponse
	for _, version := range languageVersions {
		languageVersionResponses = append(languageVersionResponses, api.LanguageVersionResponse{
			ID:             version.ID,
			LanguageID:     version.LanguageID,
			Version:        version.Version,
			NixPackageName: version.NixPackageName,
			Template:       version.Template.String,
			SearchQuery:    version.SearchQuery,
			DefaultVersion: version.DefaultVersion,
		})
	}

	return &api.GetAllLanguageVersionsOK{
		LanguageVersions: languageVersionResponses,
	}, nil
}

// GetLanguageById implements api.Handler.
func (s *OdinServer) GetLanguageById(ctx context.Context, params api.GetLanguageByIdParams) (api.GetLanguageByIdRes, error) {
	language, err := s.queries.GetLanguageByID(ctx, params.ID)
	if err != nil {
		s.logger.Printf("Error fetching language with ID %d: %v", params.ID, err)

		return &api.GetLanguageByIdInternalServerError{
			Message: "An internal error occurred while fetching the language. Please try again later.",
		}, err
	}

	if language.ID == 0 {
		return &api.GetLanguageByIdNotFound{
			Message: fmt.Sprintf("Language with ID %d not found.", params.ID),
		}, nil
	}

	response := api.LanguageResponse{
		ID:             language.ID,
		Name:           language.Name,
		Extension:      language.Extension,
		MonacoLanguage: language.MonacoLanguage,
		DefaultCode:    language.DefaultCode,
	}

	return &api.GetLanguageByIdOK{
		Language: response,
	}, nil
}

// GetLanguageVersionById implements api.Handler.
func (s *OdinServer) GetLanguageVersionById(ctx context.Context, params api.GetLanguageVersionByIdParams) (api.GetLanguageVersionByIdRes, error) {
	version, err := s.queries.GetLanguageVersionByID(ctx, params.ID)
	if err != nil {
		s.logger.Printf("Error fetching language version with ID %d: %v", params.ID, err)

		return &api.GetLanguageVersionByIdInternalServerError{
			Message: "An internal error occurred while fetching the language version. Please try again later.",
		}, err
	}

	if version.ID == 0 {
		return &api.GetLanguageVersionByIdNotFound{
			Message: fmt.Sprintf("Language version with ID %d not found.", params.ID),
		}, nil
	}

	response := api.LanguageVersionResponse{
		ID:             version.ID,
		LanguageID:     version.LanguageID,
		Version:        version.Version,
		NixPackageName: version.NixPackageName,
		Template:       version.Template.String,
		SearchQuery:    version.SearchQuery,
		DefaultVersion: version.DefaultVersion,
	}

	return &api.GetLanguageVersionByIdOK{
		Language: response,
	}, nil
}

// UpdateLanguage implements api.Handler.
func (s *OdinServer) UpdateLanguage(ctx context.Context, req *api.Language, params api.UpdateLanguageParams) (api.UpdateLanguageRes, error) {
	dbParams := db.UpdateLanguageParams{
		ID:             params.ID,
		Name:           req.Name,
		Extension:      req.Extension,
		MonacoLanguage: req.MonacoLanguage,
		DefaultCode:    req.DefaultCode,
	}

	id, err := s.queries.UpdateLanguage(ctx, dbParams)
	if err != nil {
		return &api.UpdateLanguageInternalServerError{
			Message: fmt.Sprintf("Failed to update language: %v", err),
		}, err
	}

	if id == 0 {
		return &api.UpdateLanguageBadRequest{
			Message: fmt.Sprintf("Language with ID %d not found: %v", params.ID, err),
		}, nil
	}

	response := api.LanguageResponse{
		ID:             id,
		Name:           req.Name,
		Extension:      req.Extension,
		MonacoLanguage: req.MonacoLanguage,
		DefaultCode:    req.DefaultCode,
	}

	return &api.UpdateLanguageOK{
		Language: response,
	}, nil
}

// UpdateLanguageVersion implements api.Handler.
func (s *OdinServer) UpdateLanguageVersion(ctx context.Context, req *api.LanguageVersion, params api.UpdateLanguageVersionParams) (api.UpdateLanguageVersionRes, error) {
	dbParams := db.UpdateLanguageVersionParams{
		ID:             params.ID,
		LanguageID:     req.LanguageID,
		Version:        req.Version,
		NixPackageName: req.NixPackageName,
		Template:       pgtype.Text{String: req.Template, Valid: true},
		SearchQuery:    req.SearchQuery,
		DefaultVersion: req.DefaultVersion,
	}

	id, err := s.queries.UpdateLanguageVersion(ctx, dbParams)
	if err != nil {
		return &api.UpdateLanguageVersionInternalServerError{
			Message: fmt.Sprintf("Failed to update language version: %v", err),
		}, err
	}

	if id == 0 {
		return &api.UpdateLanguageVersionBadRequest{
			Message: fmt.Sprintf("Language version with ID %d not found: %v", params.ID, err),
		}, nil
	}

	response := api.LanguageVersionResponse{
		ID:             id,
		LanguageID:     params.ID,
		Version:        req.Version,
		NixPackageName: req.NixPackageName,
		Template:       req.Template,
		SearchQuery:    req.SearchQuery,
		DefaultVersion: req.DefaultVersion,
	}

	return &api.UpdateLanguageVersionOK{
		Language: response,
	}, nil
}

// GetAllVersions implements api.Handler.
func (s *OdinServer) GetAllVersions(ctx context.Context, params api.GetAllVersionsParams) (api.GetAllVersionsRes, error) {
	languageVersions, err := s.queries.GetVersionsByLanguageID(ctx, params.ID)
	if err != nil {
		s.logger.Printf("Error fetching language versions: %v", err)

		return &api.GetAllVersionsInternalServerError{
			Message: "An internal error occurred while fetching language versions. Please try again later.",
		}, err
	}

	if len(languageVersions) == 0 {
		return &api.GetAllVersionsOK{
			LanguageVersions: []api.LanguageVersionResponse{},
		}, nil
	}

	var languageVersionResponses []api.LanguageVersionResponse
	for _, version := range languageVersions {
		languageVersionResponses = append(languageVersionResponses, api.LanguageVersionResponse{
			ID:             version.ID,
			LanguageID:     version.LanguageID,
			Version:        version.Version,
			NixPackageName: version.NixPackageName,
			Template:       version.Template.String,
			SearchQuery:    version.SearchQuery,
			DefaultVersion: version.DefaultVersion,
		})
	}

	return &api.GetAllVersionsOK{
		LanguageVersions: languageVersionResponses,
	}, nil
}
