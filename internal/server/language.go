package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/pkg/api"
)

// GetAllLanguages implements api.Handler.
func (s *ValkyrieServer) GetAllLanguages(ctx context.Context, params api.GetAllLanguagesParams) (api.GetAllLanguagesRes, error) {
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
func (s *ValkyrieServer) GetAllLanguageVersions(ctx context.Context, params api.GetAllLanguageVersionsParams) (api.GetAllLanguageVersionsRes, error) {
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
			NixPackageName: version.NixPackageName.String,
			Template:       version.Template.String,
			DefaultVersion: version.DefaultVersion,
		})
	}

	return &api.GetAllLanguageVersionsOK{
		LanguageVersions: languageVersionResponses,
	}, nil
}

// GetLanguageById implements api.Handler.
func (s *ValkyrieServer) GetLanguageById(ctx context.Context, params api.GetLanguageByIdParams) (api.GetLanguageByIdRes, error) {
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
func (s *ValkyrieServer) GetLanguageVersionById(ctx context.Context, params api.GetLanguageVersionByIdParams) (api.GetLanguageVersionByIdRes, error) {
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
		NixPackageName: version.NixPackageName.String,
		Template:       version.Template.String,
		DefaultVersion: version.DefaultVersion,
	}

	return &api.GetLanguageVersionByIdOK{
		Language: response,
	}, nil
}

// GetAllVersions implements api.Handler.
func (s *ValkyrieServer) GetAllVersions(ctx context.Context, params api.GetAllVersionsParams) (api.GetAllVersionsRes, error) {
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
			NixPackageName: version.NixPackageName.String,
			Template:       version.Template.String,
			DefaultVersion: version.DefaultVersion,
		})
	}

	return &api.GetAllVersionsOK{
		LanguageVersions: languageVersionResponses,
	}, nil
}
