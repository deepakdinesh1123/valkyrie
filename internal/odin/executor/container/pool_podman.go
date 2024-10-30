//go:build podman && !darwin

package container

import (
	"context"
	"fmt"
	"syscall"

	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/jackc/puddle/v2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func constructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()

	var cont Container

	switch envConfig.ODIN_CONTAINER_ENGINE {
	case "docker":
		return Container{}, fmt.Errorf("docker is not supported")
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

		s.ContainerBasicConfig = specgen.ContainerBasicConfig{
			Env: map[string]string{
				"NIX_CHANNELS_ENVIRONMENT": envConfig.ODIN_NIX_CHANNELS_ENVIRONMENT,
				"NIX_USER_ENVIRONMENT":     envConfig.ODIN_NIX_USER_ENVIRONMENT,
			},
		}

		s.ContainerStorageConfig.OverlayVolumes = []*specgen.OverlayVolume{
			{
				Destination: "/nix",
				Source:      envConfig.ODIN_NIX_STORE,
			},
		}

		readOnlyFileSystem := true
		readWriteTmpfs := true
		s.ContainerSecurityConfig = specgen.ContainerSecurityConfig{
			UserNS: specgen.Namespace{
				NSMode: specgen.KeepID,
				Value:  "uid=2048,gid=2048",
			},
			ReadOnlyFilesystem: &readOnlyFileSystem,
			ReadWriteTmpfs:     &readWriteTmpfs,
			CapDrop: []string{
				"CAP_DAC_OVERRIDE",
				"CAP_FOWNER",
				"CAP_FSETID",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
				"CAP_SETGID",
				"CAP_SETPCAP",
				"CAP_SETUID",
				"CAP_SYS_CHROOT",
			},
		}

		memunit := 1024 * 1024
		mem := int64(envConfig.ODIN_WORKER_MEMORY_LIMIT * int64(memunit))

		quota := int64(300000)
		burst := uint64(100000)
		period := uint64(1000000)
		realTimeRuntime := int64(500000)
		realTimePeriod := uint64(1000000)

		s.ResourceLimits = &specs.LinuxResources{
			Memory: &specs.LinuxMemory{
				Limit: &mem,
			},
			CPU: &specs.LinuxCPU{
				Quota:           &quota,
				Burst:           &burst,
				Period:          &period,
				RealtimeRuntime: &realTimeRuntime,
				RealtimePeriod:  &realTimePeriod,
			},
		}

		containerRemove := true
		s.Remove = &containerRemove

		pdContainer, err := containers.CreateWithSpec(connection, s, nil)
		cont.ID = pdContainer.ID
		if err != nil {
			return Container{}, err
		}
		err = containers.Start(connection, pdContainer.ID, nil)
		if err != nil {
			return Container{}, err
		}
		contInfo, err := containers.Inspect(connection, pdContainer.ID, nil)
		if err != nil {
			return Container{}, err
		}
		cont.Name = contInfo.Name
		cont.PID = contInfo.State.Pid
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
	for i := 0; i < int(initPoolSize); i += 1 {
		err := pool.CreateResource(ctx)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
