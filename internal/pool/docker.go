//go:build docker || all || darwin

package pool

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var getDockerClientOnce sync.Once
var dockerclient *client.Client

func GetDockerClient() *client.Client {
	getDockerClientOnce.Do(
		func() {
			c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.Fatalf("Error getting docker client: %+v", err)
				return
			}
			dockerclient = c
		},
	)
	return dockerclient
}

func DockerExecConstructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()
	var cont Container
	dClient := GetDockerClient()
	if dClient == nil {
		return Container{}, fmt.Errorf("could not get docker client")
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  true,
		Runtime:     envConfig.CONTAINER_RUNTIME,
		NetworkMode: "bridge",
	}
	createResp, err := dClient.ContainerCreate(ctx, &container.Config{
		Image:       envConfig.EXECUTION_IMAGE,
		StopTimeout: &envConfig.WORKER_TASK_TIMEOUT,
		StopSignal:  "SIGKILL",
	},
		hostConfig,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"valkyrie-network": {},
			},
		},
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

func DockerSandboxConstructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()
	var cont Container
	dClient := GetDockerClient()
	if dClient == nil {
		return Container{}, fmt.Errorf("could not get docker client")
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  true,
		Runtime:     envConfig.CONTAINER_RUNTIME,
		NetworkMode: "bridge",
	}

	containerName := namesgenerator.GetRandomName(10)

	createResp, err := dClient.ContainerCreate(ctx, &container.Config{
		Image:      envConfig.SANDBOX_IMAGE,
		StopSignal: "SIGKILL",
		ExposedPorts: nat.PortSet{
			nat.Port("9090/tcp"): {},
			nat.Port("1618/tcp"): {},
		},
		Labels: map[string]string{
			fmt.Sprintf("traefik.http.routers.%s-cs.rule", containerName):                      fmt.Sprintf("Host(`%s-cs.%s`)", containerName, envConfig.SANDBOX_HOSTNAME),
			fmt.Sprintf("traefik.http.routers.%s-cs.entrypoints", containerName):               "http",
			fmt.Sprintf("traefik.http.routers.%s-cs.service", containerName):                   fmt.Sprintf("%s-cs", containerName),
			fmt.Sprintf("traefik.http.services.%s-cs.loadbalancer.server.port", containerName): "9090",
			fmt.Sprintf("traefik.http.routers.%s-ag.rule", containerName):                      fmt.Sprintf("Host(`%s-ag.%s`)", containerName, envConfig.SANDBOX_HOSTNAME),
			fmt.Sprintf("traefik.http.routers.%s-ag.entrypoints", containerName):               "http",
			fmt.Sprintf("traefik.http.routers.%s-ag.service", containerName):                   fmt.Sprintf("%s-ag", containerName),
			fmt.Sprintf("traefik.http.services.%s-ag.loadbalancer.server.port", containerName): "1618",
			"traefik.enable": "true",
		},
	},
		hostConfig,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"valkyrie-network": {},
			},
		},
		nil,
		containerName,
	)
	if err != nil {
		fmt.Printf("error creating container %v\n", err)
		return Container{}, err
	}
	cont.ID = createResp.ID
	err = dClient.ContainerStart(ctx, createResp.ID, container.StartOptions{})
	if err != nil {
		fmt.Printf("error starting container %v\n", err)
		return Container{}, err
	}
	contInfo, err := dClient.ContainerInspect(ctx, createResp.ID)
	if err != nil {
		fmt.Printf("error inspecting container %v\n", err)
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

func DockerExecDestructor(cont Container) {
	KillDockerContainer(cont)
}

func backup(cont Container) {}

func DockerSandboxDestructor(cont Container) {
	backup(cont)
	KillDockerContainer(cont)
}
