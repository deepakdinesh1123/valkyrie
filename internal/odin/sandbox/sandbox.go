package sandbox

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

type SandboxHandler interface {
	Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBoxJob db.FetchSandboxJobTxResult)
	Cleanup(ctx context.Context) error
	StartContainerPool(ctx context.Context, envConfig *config.EnvConfig) error
	StartOdinStore(ctx context.Context, storeImage string, storeContainerName string, containerRuntime string) error
}
