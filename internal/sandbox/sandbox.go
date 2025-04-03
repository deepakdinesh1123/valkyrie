package sandbox

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
)

type SandboxHandler interface {
	Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBoxJob db.FetchSandboxJobTxResult)
	Cleanup(ctx context.Context) error
	StartContainerPool(ctx context.Context, envConfig *config.EnvConfig) error
}
