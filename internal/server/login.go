package server

import (
	"context"
	"net/url"

	"github.com/deepakdinesh1123/valkyrie/pkg/api"
)

func (s *ValkyrieServer) GetLoginConfig(ctx context.Context) (api.GetLoginConfigRes, error) {
	providers := []api.LoginConfigProvidersItem{}
	for provider, config := range s.loginService.GetOauthProviders() {
		state, err := s.loginService.GenerateState()
		if err != nil {
			return &api.GetLoginConfigInternalServerError{
				Message: err.Error(),
			}, nil
		}
		redirectUrl := config.AuthCodeURL(state)
		u, err := url.Parse(redirectUrl)
		if err != nil {
			return &api.GetLoginConfigInternalServerError{
				Message: err.Error(),
			}, nil
		}
		providers = append(providers, api.LoginConfigProvidersItem{
			Name:        api.LoginConfigProvidersItemName(provider),
			RedirectURL: *u,
		})
	}
	return &api.LoginConfig{
		Providers: providers,
	}, nil
}

func (s *ValkyrieServer) OauthCallback(ctx context.Context, params api.OauthCallbackParams) (api.OauthCallbackRes, error) {
	err := s.loginService.ValidateOAuthState(params.State.Value)
	if err != nil {
		return &api.OauthCallbackBadRequest{
			Message: err.Error(),
		}, nil
	}

	return &api.OauthCallbackOK{}, nil
}
