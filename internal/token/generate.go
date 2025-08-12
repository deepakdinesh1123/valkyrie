package token

import (
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateValkyrieToken(claims jwt.MapClaims, envConfig *config.EnvConfig) (string, error) {
	var jwtSecretKey = []byte(envConfig.ENCKEY)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
