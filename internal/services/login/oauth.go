package login

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func getRedirectUrl(env *config.EnvConfig) string {
	if env.ENVIRONMENT == "dev" {
		return fmt.Sprintf("http://%s:%s/api/login/oauth/callback", env.SERVER_HOST, env.SERVER_PORT)
	} else {
		return fmt.Sprintf("http://%s/api/login/oauth/callback", env.SERVER_HOST)
	}
}

func configureOauthProviders(envConfig *config.EnvConfig) map[string]*oauth2.Config {
	oauthProviders := make(map[string]*oauth2.Config)

	if envConfig.ENABLE_GITHUB_OAUTH {
		oauthProviders["github"] = &oauth2.Config{
			ClientID:     envConfig.GITHUB_CLIENT_ID,
			ClientSecret: envConfig.GITHUB_CLIENT_SECRET,
			Scopes:       []string{"email"},
			RedirectURL:  getRedirectUrl(envConfig),
			Endpoint:     github.Endpoint,
		}
	}

	if envConfig.ENABLE_GOOGLE_OAUTH {
		oauthProviders["google"] = &oauth2.Config{
			ClientID:     envConfig.GOOGLE_CLIENT_ID,
			ClientSecret: envConfig.GOOGLE_CLIENT_SECRET,
			Scopes:       []string{"email"},
			RedirectURL:  getRedirectUrl(envConfig),
			Endpoint:     google.Endpoint,
		}
	}
	return oauthProviders
}

func NewLoginService(queries db.Store, envConfig *config.EnvConfig, logger *zerolog.Logger) *LoginService {
	providers := configureOauthProviders(envConfig)
	return &LoginService{
		queries:        queries,
		envConfig:      envConfig,
		logger:         logger,
		oAuthProviders: providers,
	}
}

func (l *LoginService) GenerateState() (string, error) {
	// Generate random nonce
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"nonce": base64.URLEncoding.EncodeToString(nonce),
		"exp":   time.Now().Add(10 * time.Minute).Unix(), // 10 min expiry
		"iss":   "valkyrie",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.envConfig.ENCKEY))
}

func (l *LoginService) ValidateOAuthState(state string) error {
	token, err := jwt.Parse(state, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(l.envConfig.ENCKEY), nil
	})

	if err != nil {
		return fmt.Errorf("invalid state token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid state token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	// Verify issuer
	if iss, ok := claims["iss"]; !ok || iss != "valkyrie" {
		return fmt.Errorf("invalid token issuer")
	}

	if exp, ok := claims["exp"]; ok {
		if expTime, ok := exp.(float64); ok {
			if time.Now().Unix() > int64(expTime) {
				return fmt.Errorf("state token expired")
			}
		}
	}

	return nil
}

func (l *LoginService) getUserDetails(token string) (string, error) {
	return "", nil
}

func (l *LoginService) OAuthLogin(state string, code string) (string, error) {
	err := l.ValidateOAuthState(state)
	if err != nil {
		return "", err
	}
	token, err := l.getAccessToken(code)
	if err != nil {
		return "", err
	}

	l.getUserDetails(token)
	return "", nil
}

func (l *LoginService) getAccessToken(code string) (string, error) {
	return "", nil
}

func (l *LoginService) GetOauthProviders() map[string]*oauth2.Config {
	return l.oAuthProviders
}
