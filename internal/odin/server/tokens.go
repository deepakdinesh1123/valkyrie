package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/auth"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *OdinServer) GetToken(ctx context.Context, req *api.GetTokenReq) (api.GetTokenRes, error) {
	if req.Username == s.envConfig.ODIN_USER_NAME && req.Password == s.envConfig.ODIN_USER_PASS {
		token, err := auth.GenerateToken("admin", 24, s.envConfig.ODIN_SECRET_KEY)
		if err != nil {
			return &api.GetTokenInternalServerError{
				Message: fmt.Sprintf("Failed to generate token: %v", err),
			}, nil
		}
		return &api.GetTokenOK{
			Token: token,
		}, nil
	}
	return &api.GetTokenUnauthorized{
		Message: "Invalid credentials",
	}, nil
}

func (s *OdinServer) GenerateUserToken(ctx context.Context) (api.GenerateUserTokenRes, error) {
	if auth.CheckRoles(ctx, []string{"admin"}) {
		token, err := auth.GenerateToken("user", 24, s.envConfig.ODIN_SECRET_KEY)
		if err != nil {
			return &api.GenerateUserTokenInternalServerError{
				Message: fmt.Sprintf("Failed to generate token: %v", err),
			}, nil
		}
		return &api.GenerateUserTokenOK{
			Token: token,
		}, nil
	}
	return &api.GenerateUserTokenForbidden{
		Message: "Forbidden",
	}, nil
}
