package provider

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/models"
)

type Provider interface {
	Execute(ctx context.Context, execReq db.Jobqueue) (models.ExecutionResult, error)
}
