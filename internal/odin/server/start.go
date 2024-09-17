package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/middleware"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/sync/errgroup"
)

var shutdownTimeout = time.Second * 5

func (s *OdinServer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := fmt.Sprintf("%s:%s", s.envConfig.ODIN_SERVER_HOST, s.envConfig.ODIN_SERVER_PORT)
	done := make(chan bool, 1)

	var server *http.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/executions/{executionId}/sse", s.ExecuteSSE)
	mux.HandleFunc("/executions/execute/ws", s.ExecuteWS)
	mux.Handle("/", s.server)

	route_finder := middleware.MakeRouteFinder(s.server)
	server = &http.Server{
		ReadHeaderTimeout: time.Second * 5,
		Addr:              addr,
		Handler: middleware.Wrap(mux,
			middleware.Instrument("server", route_finder, s.tp, s.mp, s.prop),
			middleware.Labeler(route_finder),
		),
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-ctx.Done()

		s.logger.Info().Msg("Shutting down server")

		var err error
		err = errors.Join(err, s.otelShutdown(context.Background()))
		if err != nil {
			s.logger.Err(err).Msg("Failed to shutdown OpenTelemetry")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		s.logger.Info().Msg(fmt.Sprintf("Starting server on %s", addr))
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				s.logger.Info().Msg("Server closed")
				return nil
			}
			s.logger.Err(err).Msg("Failed to start server")
			done <- true
		}
		return err
	})

	g.Go(func() error {
		ticker := time.NewTicker(time.Duration(s.envConfig.ODIN_JOB_PRUNE_FREQ) * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				s.logger.Info().Msg("Pruning completed jobs")
				err := s.queries.PruneCompletedJobs(ctx)
				if err != nil {
					s.logger.Err(err).Msg("Failed to prune completed jobs")
				}
				return nil
			}
		}
	})

	g.Go(func() error {
		ticker := time.NewTicker(time.Duration(10) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				ids, err := s.queries.GetStaleWorkers(ctx)
				if err != nil {
					s.logger.Err(err).Msg("Failed to get stale workers")
				}
				for _, id := range ids {
					s.logger.Info().Msg(fmt.Sprintf("Requeuing jobs for stale worker %d", id))
					err := s.queries.RequeueWorkerJobs(ctx, pgtype.Int4{Int32: id, Valid: true})
					if err != nil {
						s.logger.Err(err).Msg("Failed to requeue jobs for stale worker")
					}
				}
				err = s.queries.RequeueLTJobs(ctx)
				if err != nil {
					s.logger.Err(err).Msg("Failed to requeue jobs")
				}
			}
		}
	})
	err := g.Wait()

	if err != nil {
		s.logger.Err(err).Msg("Failed to start server")
		done <- true
	}
}
