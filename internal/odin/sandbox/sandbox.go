package sandbox

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

type SandboxHandler interface {
	Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBox db.Sandbox, logger zerolog.Logger)
}

func NewSandboxHandler() {

}
