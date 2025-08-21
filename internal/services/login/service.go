package login

import (
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type LoginService struct {
	oAuthProviders map[string]*oauth2.Config
	queries        db.Store
	envConfig      *config.EnvConfig
	logger         *zerolog.Logger
}
