package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/middleware"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (s *OdinServer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := fmt.Sprintf("%s:%s", s.envConfig.ODIN_SERVER_HOST, s.envConfig.ODIN_SERVER_PORT)
	done := make(chan bool, 1)

	var server *http.Server
	mux := http.NewServeMux()

	mux.HandleFunc("/executions/{executionId}/sse", s.ExecuteSSE)
	mux.HandleFunc("/executions/execute/ws", s.ExecuteWS)
	mux.Handle("/", s.server)

	if s.envConfig.ODIN_ENABLE_TELEMETRY {
		server = &http.Server{
			Addr:    addr,
			Handler: otelhttp.NewHandler(middleware.LoggingMiddleware(mux), "/"),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}
	} else {
		server = &http.Server{
			Addr:    addr,
			Handler: middleware.LoggingMiddleware(mux),
		}
	}

	go func() {
		s.logger.Info().Msg(fmt.Sprintf("Starting server on %s", addr))
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				s.logger.Info().Msg("Server closed")
				return
			}
			s.logger.Err(err).Msg("Failed to start server")
			done <- true
		}
	}()

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(s.envConfig.ODIN_JOB_PRUNE_FREQ) * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.logger.Info().Msg("Pruning completed jobs")
				err := s.queries.PruneCompletedJobs(ctx)
				if err != nil {
					s.logger.Err(err).Msg("Failed to prune completed jobs")
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(5) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ids, err := s.queries.GetStaleWorkers(ctx)
				if err != nil {
					s.logger.Err(err).Msg("Failed to get stale workers")
				}
				for _, id := range ids {
					s.logger.Info().Msg(fmt.Sprintf("Requeuing jobs for stale worker %d", id))
					s.queries.RequeueWorkerJobs(ctx, pgtype.Int4{Int32: id, Valid: true})
				}
			}
		}
	}(ctx)

	go func() {
		<-ctx.Done()

		s.logger.Info().Msg("Shutting down OpenTelemetry")
		var err error
		err = errors.Join(err, s.otelShutdown(context.Background()))
		if err != nil {
			s.logger.Err(err).Msg("Failed to shutdown OpenTelemetry")
		}

		s.logger.Info().Msg("Shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Err(err).Msg("Failed to shutdown server")
		}
		done <- true
	}()
	<-done
}
