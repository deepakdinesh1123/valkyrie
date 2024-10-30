//go:build docker

package container

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/puddle/v2"
)

func constructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()

	var cont Container

	switch envConfig.ODIN_CONTAINER_ENGINE {
	case "docker":
		dClient := getDockerClient()
		if dClient == nil {
			return Container{}, fmt.Errorf("could not get docker client")
		}
		createResp, err := dClient.ContainerCreate(ctx, &container.Config{
			Image:       envConfig.ODIN_WORKER_DOCKER_IMAGE,
			StopTimeout: &envConfig.ODIN_WORKER_TASK_TIMEOUT,
			StopSignal:  "SIGKILL",
		},
			&container.HostConfig{
				AutoRemove: true,
			},
			nil,
			nil,
			"",
		)
		if err != nil {
			return Container{}, err
		}
		cont.ID = createResp.ID
		err = dClient.ContainerStart(ctx, createResp.ID, container.StartOptions{})
		if err != nil {
			return Container{}, err
		}
		contInfo, err := dClient.ContainerInspect(ctx, createResp.ID)
		if err != nil {
			return Container{}, err
		}
		cont.Name = contInfo.Name
		cont.PID = contInfo.State.Pid
	case "podman":
		return cont, fmt.Errorf("Podman is not supported")
	}

	return cont, nil
}

func destructor(cont Container) {
	fmt.Println("killing container", cont.Name)
	KillContainer(cont.PID)
}

func NewContainerPool(ctx context.Context, initPoolSize int32, maxPoolSize int32) (*puddle.Pool[Container], error) {
	pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: constructor, Destructor: destructor, MaxSize: maxPoolSize})
	if err != nil {
		return nil, err
	}
	fmt.Println("The initial container pool size is ", initPoolSize)
	for i := 0; i < int(initPoolSize); i += 1 {
		err := pool.CreateResource(ctx)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
