//go:build docker || all

package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var getDockerClientOnce sync.Once
var dockerclient *client.Client

func GetDockerClient() *client.Client {
	getDockerClientOnce.Do(
		func() {
			c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return
			}
			dockerclient = c
		},
	)
	return dockerclient
}

func DockerConstructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()
	var cont Container
	dClient := GetDockerClient()
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
	cont.Engine = "docker"
	return cont, nil
}

func KillDockerContainer(cont Container) error {
	client := GetDockerClient()
	return client.ContainerRemove(context.TODO(), cont.ID, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
}

func DockerDestructor(cont Container) {
	KillDockerContainer(cont)
}
