package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/go-chi/jwtauth/v5"
)

type SecurityHandlerImpl struct {
	envConfig *config.EnvConfig
	ja        *jwtauth.JWTAuth
}

// NewSecurityHandler creates a new instance of SecurityHandlerImpl
func NewSecurityHandler(envConfig *config.EnvConfig, ja *jwtauth.JWTAuth) *SecurityHandlerImpl {
	return &SecurityHandlerImpl{
		envConfig: envConfig,
		ja:        ja,
	}
}

func (s *SecurityHandlerImpl) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	token, err := s.ja.Decode(t.Token)
	if err != nil {
		return ctx, fmt.Errorf("invalid token: %w", err)
	}
	if token == nil {
		return ctx, fmt.Errorf("invalid token")
	}
	return ctx, nil
}

func (s *SecurityHandlerImpl) HandleXAuthToken(ctx context.Context, operationName api.OperationName, t api.XAuthToken) (context.Context, error) {
	switch t.APIKey {
	case s.envConfig.ADMIN_TOKEN:
		ctx = context.WithValue(ctx, config.RoleKey, "admin")
	default:
		ctx = context.WithValue(ctx, config.RoleKey, "user")
	}
	return ctx, nil
}
