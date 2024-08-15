package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/middleware"
)

func (s *OdinServer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := fmt.Sprintf("%s:%s", s.envConfig.ODIN_SERVER_HOST, s.envConfig.ODIN_SERVER_PORT)
	done := make(chan bool, 1)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: middleware.LoggingMiddleware(s.server),
	}

	go func() {
		s.logger.Info().Msg(fmt.Sprintf("Starting server on %s", addr))
		err := httpServer.ListenAndServe()
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

	go func() {
		<-ctx.Done()
		s.logger.Info().Msg("Shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Err(err).Msg("Failed to shutdown server")
		}
		done <- true
	}()
	<-done
}
