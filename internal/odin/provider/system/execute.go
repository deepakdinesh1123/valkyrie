package system

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/models"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

func (s *SystemProvider) Execute(ctx context.Context, execReq db.Jobqueue) (models.ExecutionResult, error) {
	return models.ExecutionResult{}, nil
}
