package provider

import (
	"context"
	"os"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/docker"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/system"
	"github.com/rs/zerolog"
)

type Provider interface {
	Execute(ctx context.Context, wg *sync.WaitGroup, execReq db.Jobqueue) error
}

func GetProvider(ctx context.Context, queries *db.Queries, envConfig *config.EnvConfig, logger *zerolog.Logger) (Provider, error) {
	var provider Provider
	var err error
	switch envConfig.ODIN_WORKER_PROVIDER {
	case "docker":
		provider, err = docker.NewDockerProvider()
		if err != nil {
			logger.Err(err).Msg("Failed to create docker provider")
			return nil, err
		}
	case "system":
		if _, err := os.Stat(envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR); os.IsNotExist(err) {
			err = os.Mkdir(envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, os.ModePerm)
			if err != nil {
				logger.Err(err).Msg("Failed to create system provider base directory")
				return nil, err
			}
		}
		provider, err = system.NewSystemProvider(envConfig, queries, logger)
		if err != nil {
			logger.Err(err).Msg("Failed to create system provider")
			return nil, err
		}
	default:
		logger.Err(err).Msg("Invalid provider")
		return nil, err
	}
	return provider, nil
}
