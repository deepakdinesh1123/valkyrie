package auth

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}
type CtxStrKey string

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a *AuthHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	envconfig, err := config.GetEnvConfig()
	if err != nil {
		return nil, err
	}
	_, claims, err := VerifyToken(t.Token, envconfig.ODIN_SECRET_KEY)
	if err != nil {
		return nil, err
	}
	roles, err := claims.GetAudience()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, CtxStrKey("Roles"), roles)
	return ctx, nil
}

func CheckRoles(ctx context.Context, roles []string) bool {

	ctx_role_set := make(map[string]string)
	roles_in_ctx := ctx.Value(CtxStrKey("Roles")).(jwt.ClaimStrings)

	for _, role := range roles {
		ctx_role_set[role] = role
	}

	for _, role := range roles_in_ctx {
		fmt.Println(role)
		if _, ok := ctx_role_set[role]; ok {
			return true
		}
	}
	return false
}
