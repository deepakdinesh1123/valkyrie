package container

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
)

type ContainerClient interface {
	WriteFiles(ctx context.Context, containerID string, prepDir string, job *db.Job) error
	GetContainer(ctx context.Context, execReq db.ExecRequest) (string, error)
	CheckImageExists(ctx context.Context, imageName string) error
	BuildImage(ctx context.Context, imageName string) error
	Execute(ctx context.Context, containerID string, command []string) (bool, string, error)
	DestroyContainer(ctx context.Context, containerId string)
	Cleanup(ctx context.Context)
}
