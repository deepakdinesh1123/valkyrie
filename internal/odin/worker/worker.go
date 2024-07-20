package worker

import (
	"context"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type Worker struct {
	queries  *db.Queries
	env      *config.EnvConfig
	provider provider.Provider
	logger   *zerolog.Logger
}

func NewWorker(ctx context.Context, queries *db.Queries, env *config.EnvConfig, prvdr provider.Provider, logger *zerolog.Logger) *Worker {
	return &Worker{
		queries:  queries,
		env:      env,
		provider: prvdr,
		logger:   logger,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Duration(w.env.ODIN_WORKER_POLL_FREQ) * time.Second)
			job, err := w.queries.FetchJob(ctx)
			if err != nil {
				if err == pgx.ErrNoRows {
					continue
				}
				w.logger.Err(err).Msgf("Worker: failed to fetch job")
				return err
			}
			w.logger.Info().Msgf("Worker: fetched job %d", job.ID)
			_, err = w.provider.Execute(ctx, job)
			w.logger.Info().Msgf("Worker: executed job %d", job.ID)
			if err != nil {
				return err
			}
		}
	}
}
