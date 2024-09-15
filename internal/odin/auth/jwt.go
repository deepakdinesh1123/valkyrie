package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(role string, expiration int, secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "odin",
		"aud": role,
		"exp": time.Now().Add(time.Duration(expiration) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	return claims.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string, secretKey string) (*jwt.Token, *jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}
	return token, &claims, nil
}
