package system

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
)

func (s *SystemProvider) Execute(ctx context.Context, execReq db.Jobqueue) (models.ExecutionResult, error) {
	go func() {
		s.logger.Info().Msg("System: executing job")
		<-ctx.Done()
		s.logger.Info().Msg("System: Exiting routine")
	}()
	return models.ExecutionResult{}, nil
}
