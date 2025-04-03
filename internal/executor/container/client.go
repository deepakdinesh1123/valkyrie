package container

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/jackc/puddle/v2"
)

type ContainerClient interface {
	WriteFiles(ctx context.Context, containerID string, prepDir string, job *db.Job) error
	GetContainer(ctx context.Context) (*puddle.Resource[pool.Container], error)
	Execute(ctx context.Context, containerID string, command []string) (bool, string, error)
}
