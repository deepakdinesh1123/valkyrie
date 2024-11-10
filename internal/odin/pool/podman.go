//go:build (podman || all) && !darwin

package pool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var getPodmanConnectionOnce sync.Once
var podmanConnection context.Context

func GetPodmanConnection() context.Context {
	getPodmanConnectionOnce.Do(func() {
		sock_dir := os.Getenv("XDG_RUNTIME_DIR")
		socket := "unix:" + sock_dir + "/podman/podman.sock"
		pc, err := bindings.NewConnection(context.Background(), socket)
		if err != nil {
			return
		}
		podmanConnection = pc
	})
	return podmanConnection
}

func PodConstructor(ctx context.Context) (Container, error) {
	envConfig, _ := config.GetEnvConfig()

	var cont Container
	connection := GetPodmanConnection()
	if connection == nil {
		return Container{}, fmt.Errorf("could not get podman connection")
	}

	// _, err := volumes.Create(connection, types.VolumeCreateOptions{
	// 	Name:           "shared-cache",
	// 	IgnoreIfExists: true,
	// }, &volumes.CreateOptions{})

	// if err != nil {
	// 	return Container{}, fmt.Errorf("could not create volume: %s", err)
	// }
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

	// s.ContainerStorageConfig.Volumes = []*specgen.NamedVolume{
	// 	{
	// 		Dest: "/home/valnix/.cache/cached-nix-shell",
	// 		Name: "shared-cache",
	// 	},
	// }
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
		return Container{}, fmt.Errorf("could not create container: %s", err)
	}
	err = containers.Start(connection, pdContainer.ID, nil)
	if err != nil {
		return Container{}, fmt.Errorf("could not start container: %s", err)
	}
	contInfo, err := containers.Inspect(connection, pdContainer.ID, nil)
	if err != nil {
		return Container{}, fmt.Errorf("could not inspect container: %s", err)
	}
	cont.Name = contInfo.Name
	cont.PID = contInfo.State.Pid
	cont.Engine = "podman"

	return cont, nil
}

func Poddestructor(cont Container) {
	KillContainer(cont.PID)
}

func KillContainer(pid int) error {
	_, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Container with given PID has already been killed")
	}

	cmd := exec.Command("kill", "-KILL", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill container: %w", err)
	}
	return nil
}