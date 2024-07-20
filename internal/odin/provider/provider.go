package provider

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/models"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

type Provider interface {
	Execute(ctx context.Context, execReq db.Jobqueue) (models.ExecutionResult, error)
}
