package docker

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
)

func (d *DockerProvider) Execute(ctx context.Context, execReq db.Jobqueue) (models.ExecutionResult, error) {
	return models.ExecutionResult{}, nil
}
