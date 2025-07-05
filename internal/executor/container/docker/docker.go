//go:build docker || all || darwin

package docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/executor/container/common"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/deepakdinesh1123/valkyrie/internal/secret"
	"github.com/deepakdinesh1123/valkyrie/internal/services/execution"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DockerProvider struct {
	client    *client.Client
	queries   db.Store
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
}

func GetDockerProvider(envConfig *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*DockerProvider, error) {
	client := pool.GetDockerClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &DockerProvider{
		client:    client,
		queries:   queries,
		envConfig: envConfig,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}

func (d *DockerProvider) WriteFiles(ctx context.Context, containerID string, prepDir string, job *db.Job) error {

	execReq, err := d.queries.GetExecRequest(ctx, job.Arguments.ExecConfig.ExecReqId)
	if err != nil {
		return err
	}
	script, spec, err := execution.ConvertExecSpecToScript(ctx, &execReq, d.queries)
	if err != nil {
		return fmt.Errorf("error writing files: %s", err)
	}
	files := map[string]string{
		"exec.sh":       script,
		spec.ScriptName: execReq.Code.String,
		"input.txt":     execReq.Input.String,
	}

	tarFilePath, err := common.CreateTarArchive(files, execReq.Files, prepDir)
	if err != nil {
		return err
	}
	defer os.Remove(tarFilePath)

	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		d.logger.Err(err).Msg("Failed to open tar file")
		return err
	}
	defer tarFile.Close()
	err = d.client.CopyToContainer(
		ctx,
		containerID,
		filepath.Join("/home/valnix/valkyrie"),
		tarFile,
		container.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: true},
	)
	if err != nil {
		d.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	return nil
}

func (d *DockerProvider) GetContainer(ctx context.Context, execReq db.ExecRequest) (string, error) {
	imageName := fmt.Sprintf("valkyrie/shell/%s", strings.Join(execReq.SystemDependencies, "/"))

	err := d.CheckImageExists(ctx, imageName)
	if err != nil {
		if client.IsErrNotFound(err) {
			err = d.BuildImage(ctx, imageName)
			if err != nil {
				return "", fmt.Errorf("error building image: %v", err)
			}
		} else {
			return "", err
		}
	}

	containerConfig := container.Config{
		Image:       imageName,
		StopTimeout: &d.envConfig.WORKER_MAX_TASK_TIMEOUT,
		StopSignal:  "SIGKILL",
		Labels: map[string]string{
			"valkyrie": "execution",
		},
	}

	if len(execReq.Secrets) > 0 {
		secretsMap, err := secret.DecodeSecrets(execReq.Secrets, d.envConfig.ENCKEY)
		if err != nil {
			return "", err
		}
		containerConfig.Env = common.ConvertSecretsMapToSlice(secretsMap)
	}

	containerCreateResp, err := d.client.ContainerCreate(ctx, &containerConfig,
		&container.HostConfig{
			AutoRemove:  true,
			Runtime:     d.envConfig.CONTAINER_RUNTIME,
			NetworkMode: "bridge",
		}, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"devpi-network": {},
			},
		}, nil, "")
	if err != nil {
		return "", fmt.Errorf("error creating container: %v", err)
	}

	err = d.client.ContainerStart(ctx, containerCreateResp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("error starting container: %v", err)
	}

	contInfo, err := d.client.ContainerInspect(ctx, containerCreateResp.ID)
	if err != nil {
		return "", fmt.Errorf("error inspecting container: %v", err)
	}

	return contInfo.ID, nil
}

func (d *DockerProvider) CheckImageExists(ctx context.Context, imageName string) error {
	_, _, err := d.client.ImageInspectWithRaw(ctx, imageName)
	return err
}

func (d *DockerProvider) BuildImage(ctx context.Context, imageName string) error {
	// Read the Dockerfile template
	dockerfileTmpl, err := execution.ExecTemplates.ReadFile("templates/containerImage.tmpl")
	if err != nil {
		return fmt.Errorf("error reading image template: %v", err)
	}
	d.logger.Info().Str("template", string(dockerfileTmpl)).Msg("docker template loaded")

	nixRunSh, err := execution.ExecScripts.ReadFile("scripts/nix_run.sh")
	if err != nil {
		return fmt.Errorf("error reading nix run script: %v", err)
	}

	nixStopSh, err := execution.ExecScripts.ReadFile("scripts/nix_stop.sh")
	if err != nil {
		return fmt.Errorf("error reading nix stop script: %v", err)
	}

	// Template arguments structure
	type ContainerImageTmplArgs struct {
		NIXERY_IMAGE string
		BASE_IMAGE   string
	}

	// Execute template to generate Dockerfile content
	var dockerImg bytes.Buffer
	dfTemplate, err := template.New("dockerfile").Parse(string(dockerfileTmpl))
	if err != nil {
		return fmt.Errorf("error parsing dockerfile template: %v", err)
	}

	err = dfTemplate.Execute(&dockerImg, ContainerImageTmplArgs{
		NIXERY_IMAGE: strings.ReplaceAll(imageName, "valkyrie", d.envConfig.NIXERY_URL),
		BASE_IMAGE:   d.envConfig.EXECUTION_IMAGE,
	})
	if err != nil {
		return fmt.Errorf("error executing dockerfile template: %v", err)
	}
	d.logger.Info().Str("dockerfile", dockerImg.String()).Msg("generated dockerfile")

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "docker-build-*")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Write Dockerfile to temp directory
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, dockerImg.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing dockerfile to temp dir: %v", err)
	}

	// Create scripts directory in temp dir
	scriptsDir := filepath.Join(tempDir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return fmt.Errorf("error creating scripts directory: %v", err)
	}

	// Write scripts to temp directory
	nixRunPath := filepath.Join(scriptsDir, "nix_run.sh")
	if err := os.WriteFile(nixRunPath, nixRunSh, 0755); err != nil {
		return fmt.Errorf("error writing nix_run.sh to temp dir: %v", err)
	}

	nixStopPath := filepath.Join(scriptsDir, "nix_stop.sh")
	if err := os.WriteFile(nixStopPath, nixStopSh, 0755); err != nil {
		return fmt.Errorf("error writing nix_stop.sh to temp dir: %v", err)
	}

	// Create tar archive from temp directory
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Walk through temp directory and add all files to tar
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == tempDir {
			return nil
		}

		// Get relative path from temp directory
		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If it's a regular file, write its content
		if info.Mode().IsRegular() {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if _, err := tw.Write(fileContent); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		tw.Close()
		return fmt.Errorf("error creating tar archive: %v", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("error closing tar writer: %v", err)
	}

	// Build Docker image
	dockerfileTarReader := bytes.NewReader(buf.Bytes())
	buildOptions := types.ImageBuildOptions{
		Tags:           []string{imageName},
		Dockerfile:     "Dockerfile",
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		NoCache:        false,
		SuppressOutput: false,
	}

	d.logger.Info().Str("imageName", imageName).Msg("starting docker build")
	resp, err := d.client.ImageBuild(ctx, dockerfileTarReader, buildOptions)
	if err != nil {
		return fmt.Errorf("error building docker image: %v", err)
	}
	defer resp.Body.Close()

	// Process build output and check for errors
	var buildError error
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		d.logger.Debug().Str("build_output", line).Msg("docker build")

		// Parse JSON output to check for errors
		var buildMsg map[string]any
		if err := json.Unmarshal([]byte(line), &buildMsg); err == nil {
			if errorDetail, exists := buildMsg["errorDetail"]; exists {
				buildError = fmt.Errorf("docker build failed: %v", errorDetail)
				break
			}
			if stream, exists := buildMsg["stream"]; exists {
				d.logger.Info().Str("stream", fmt.Sprintf("%v", stream)).Msg("docker build progress")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading build response: %v", err)
	}

	if buildError != nil {
		return buildError
	}

	d.logger.Info().Str("imageName", imageName).Msg("docker image built successfully")
	return nil
}

func (d *DockerProvider) Execute(ctx context.Context, containerID string, command []string) (bool, string, error) {
	done := make(chan bool)

	var dexec types.IDResponse
	var err error

	go func() {
		defer func() {
			done <- true
		}()

		dexec, err = d.client.ContainerExecCreate(
			ctx,
			containerID,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				Cmd:          command,
			},
		)
		if err != nil {
			return
		}
		err = d.client.ContainerExecStart(ctx, dexec.ID, container.ExecAttachOptions{})
		if err != nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				d.logger.Info().Msg("Timelimit exceeded")
				return
			default:
				execInfo, err := d.client.ContainerExecInspect(ctx, dexec.ID)
				if err != nil {
					return
				}
				if !execInfo.Running {
					d.logger.Info().Int("Exit Code", execInfo.ExitCode).Msg("Execution process exit")
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			if dexec.ID != "" {
				stopExec, err := d.client.ContainerExecCreate(
					context.TODO(),
					containerID,
					container.ExecOptions{
						Cmd: []string{"sh", "nix_stop.sh"},
					},
				)
				if err != nil {
					return false, "", fmt.Errorf("could not create exec for nix stop: %s", err)
				}
				err = d.client.ContainerExecStart(context.TODO(), stopExec.ID, container.ExecStartOptions{})
				if err != nil {
					return false, "", fmt.Errorf("could not start the nix_stop script: %s", err)
				}
			}
			out, err := d.ReadExecLogs(context.TODO(), containerID)
			if err != nil {
				return false, "", fmt.Errorf("error reading output: %s", err)
			}
			return true, out, nil
		case context.Canceled:
			return false, "", fmt.Errorf("context canceled")
		}
	case <-done:
		out, err := d.ReadExecLogs(context.TODO(), containerID)
		if err != nil {
			return false, "", fmt.Errorf("error reading output: %s", err)
		}
		return true, out, nil
	}
	return false, "", nil
}

func (d *DockerProvider) ReadExecLogs(ctx context.Context, containerID string) (string, error) {
	var out []byte
	dexec, err := d.client.ContainerExecCreate(
		ctx,
		containerID,
		container.ExecOptions{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"sh", "-c", "cat ~/valkyrie/output.txt"},
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not create exec: %s", err)
	}
	resp, err := d.client.ContainerExecAttach(ctx, dexec.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("could not attach to container: %s", err)
	}
	if resp.Reader != nil {
		out, err = io.ReadAll(resp.Reader)
		if err != nil {
			return "", fmt.Errorf("could not read from hijacked response: %s", err)
		}
	}
	return stripCtlAndExtFromUTF8(string(out)), nil
}

func (d *DockerProvider) DestroyContainer(ctx context.Context, containerId string) {
	err := d.client.ContainerRemove(ctx, containerId, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		fmt.Printf("error destroying container %v\n", err)
	}
}

func (d *DockerProvider) Cleanup(ctx context.Context) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("valkyrie", "execution")

	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		Filters: filterArgs,
	})
	if err != nil {

	}
	for _, container := range containers {
		d.DestroyContainer(ctx, container.ID)
	}
}

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 || r == 10 {
			return r
		}
		return -1
	}, str)
}
