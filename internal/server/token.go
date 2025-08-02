package server

import (
	"context"
	"fmt"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/token"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/golang-jwt/jwt/v5"
)

func (s *ValkyrieServer) GetValkyrieToken(ctx context.Context) (api.GetValkyrieTokenRes, error) {
	role, _ := ctx.Value(config.RoleKey).(string)

	var expiryDuration time.Duration
	if role == "admin" {
		expiryDuration = time.Hour * 24 * 90
	} else {
		expiryDuration = time.Hour * 2
	}

	expiresAt := time.Now().Add(expiryDuration)

	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}

	tkn, err := token.GenerateValkyrieToken(claims, s.envConfig)
	if err != nil {
		return &api.GetValkyrieTokenInternalServerError{
			Message: fmt.Sprintf("failed to generate token: %v", err),
		}, nil
	}

	return &api.GetValkyrieTokenOK{
		Token:   tkn,
		Expires: expiresAt,
	}, nil
}
