package container

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
)

type ContainerClient interface {
	WriteFiles(ctx context.Context, containerID string, prepDir string, job db.Job) error
	GetContainer(ctx context.Context, prepDir string) (Container, error)
	Execute(ctx context.Context, containerID string, command []string) error
	ReadOutput(ctx context.Context, containerID string) (string, error)
}
