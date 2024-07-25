package worker

import (
	"context"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type Worker struct {
	ID       int
	Name     string
	queries  *db.Queries
	env      *config.EnvConfig
	provider provider.Provider
	logger   *zerolog.Logger
}

func GetWorker(ctx context.Context, name string, queries *db.Queries, env *config.EnvConfig, prvdr provider.Provider, logger *zerolog.Logger) (*Worker, error) {
	wrkr, err := queries.GetWorker(ctx, pgtype.Text{String: name, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			wrkr, err = queries.InsertWorker(ctx, pgtype.Text{String: name, Valid: true})
			if err != nil {
				logger.Err(err).Msg("Worker: failed to insert worker")
				return nil, err
			}
		} else {
			logger.Err(err).Msg("Worker: failed to get worker")
			return nil, err
		}
	}
	return &Worker{
		queries:  queries,
		env:      env,
		provider: prvdr,
		logger:   logger,
		Name:     name,
		ID:       int(wrkr.ID),
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(w.env.ODIN_WORKER_POLL_FREQ) * time.Second)
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			ticker.Stop()
			switch err {
			case context.Canceled:
				w.logger.Info().Msg("Worker: context canceled")
				return nil
			default:
				w.logger.Err(err).Msg("Worker: context error")
				return err
			}
		case <-ticker.C:
			job, err := w.queries.FetchJob(ctx)
			if err != nil {
				switch err {
				case pgx.ErrNoRows:
					continue
				case context.Canceled:
					w.logger.Info().Msg("Worker: context canceled")
					return nil
				default:
					w.logger.Err(err).Msgf("Worker: failed to fetch job")
					return err
				}
			}
			w.logger.Info().Msgf("Worker: fetched job %d", job.ID)
			_, err = w.provider.Execute(ctx, job)
			if err != nil {
				return err
			}
		}
	}
}
