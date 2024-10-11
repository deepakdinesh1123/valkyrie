package container

import (
	"context"
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/puddle/v2"
	spec "github.com/opencontainers/runtime-spec/specs-go"
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
				Binds: []string{
					fmt.Sprintf("%s:/home/odnix/.nix-portable", filepath.Join(envConfig.USER_HOME_DIR, ".nix-portable")),
				},
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
		connection := getPodmanConnection()
		if connection == nil {
			return Container{}, fmt.Errorf("could not get podman connection")
		}
		s := specgen.NewSpecGenerator(
			envConfig.ODIN_WORKER_PODMAN_IMAGE,
			false,
		)
		stopTimeout := uint(envConfig.ODIN_WORKER_TASK_TIMEOUT)
		s.StopTimeout = &stopTimeout
		stopSignal := syscall.SIGKILL
		s.StopSignal = &stopSignal
		s.OCIRuntime = "crun"

		s.ContainerStorageConfig.Mounts = append(s.ContainerStorageConfig.Mounts, spec.Mount{
			Destination: "/home/odnix/.nix-portable",
			Type:        "bind",
			Source:      filepath.Join(envConfig.USER_HOME_DIR, "merged"),
			Options:     []string{"U"},
		})

		containerRemove := true
		s.Remove = &containerRemove

		pdContainer, err := containers.CreateWithSpec(connection, s, nil)
		cont.ID = pdContainer.ID
		if err != nil {
			return Container{}, err
		}
		fmt.Println("Container created")
		err = containers.Start(connection, pdContainer.ID, nil)
		if err != nil {
			return Container{}, err
		}
		fmt.Println("Container started")
		contInfo, err := containers.Inspect(connection, pdContainer.ID, nil)
		if err != nil {
			return Container{}, err
		}
		fmt.Println("Container inspected")
		cont.Name = contInfo.Name
		cont.PID = contInfo.State.Pid
	}

	return cont, nil
}

func destructor(cont Container) {
	Cleanup(cont.HostPrepDir)
	KillContainer(cont.PID)
}

func NewContainerPool(ctx context.Context, initPoolSize int32, maxPoolSize int32) (*puddle.Pool[Container], error) {
	pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: constructor, Destructor: destructor, MaxSize: maxPoolSize})
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(initPoolSize); i += 1 {
		fmt.Println("Creating source")
		err := pool.CreateResource(ctx)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
