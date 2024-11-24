package sandbox

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

type SandboxHandler interface {
	Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBox db.Sandbox)
}
