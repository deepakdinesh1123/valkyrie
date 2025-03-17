package executor

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/rs/zerolog"
)

type Executor interface {
	Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq *db.Job, logger zerolog.Logger)
	Cleanup()
}
