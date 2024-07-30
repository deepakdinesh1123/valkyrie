package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/text/encoding/charmap"
)

func (d *DockerProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Jobqueue) {
	defer wg.Done()
	containerName := namesgenerator.GetRandomName(0)
	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:       "alpinix",
			WorkingDir:  d.envConfig.ODIN_DOCKER_PROVIDER_BASE_DIR,
			StopTimeout: &d.envConfig.ODIN_WORKER_TASK_TIMEOUT,
			StopSignal:  "SIGINT",
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		d.logger.Err(err).Msg("Failed to create container")
		d.updateJob(ctx, execReq.ID, err.Error())
		return
	}
	err = d.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		d.logger.Err(err).Msg("Failed to start container")
		d.updateJob(ctx, execReq.ID, err.Error())
		return
	}

	if _, err := d.client.ContainerInspect(ctx, resp.ID); err != nil {
		d.logger.Err(err).Msg("Failed to inspect container")
		d.updateJob(ctx, execReq.ID, err.Error())
		return
	}

	err = d.writeFiles(ctx, containerName, d.envConfig.ODIN_DOCKER_PROVIDER_BASE_DIR, execReq)
	if err != nil {
		d.logger.Err(err).Msg("Failed to write files")
		d.client.ContainerKill(ctx, containerName, "SIGKILL")
		d.updateJob(ctx, execReq.ID, err.Error())
		return
	}

	done := make(chan bool, 1)

	go func() {
		resp, err := d.client.ContainerExecCreate(
			ctx,
			containerName,
			container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				WorkingDir:   d.envConfig.ODIN_DOCKER_PROVIDER_BASE_DIR,
				Cmd:          []string{"nix", "run"},
			},
		)
		if err != nil {
			d.logger.Err(err).Msg("Failed to create exec")
			d.updateJob(ctx, execReq.ID, err.Error())
			done <- true
			return
		}
		hijResp, err := d.client.ContainerExecAttach(ctx, resp.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msg("Failed to attach exec")
			d.updateJob(ctx, execReq.ID, err.Error())
			done <- true
			return
		}
		var out []byte
		for {
			out = make([]byte, 1024)
			_, err := hijResp.Reader.Read(out)
			if err != nil {
				d.logger.Err(err).Msg("Failed to read output")
				break
			}
		}
		decodedOut, err := charmap.ISO8859_1.NewDecoder().Bytes(out)
		if err != nil {
			d.logger.Err(err).Msg("Failed to decode output")
			d.updateJob(ctx, execReq.ID, err.Error())
			<-done
			return
		}
		d.logger.Info().Bytes("Output", decodedOut).Msg("Writing output to stdout")
		err = d.updateJob(ctx, execReq.ID, string(decodedOut))
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		done <- true
		defer hijResp.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				d.logger.Info().Msg("Context canceled wating for process to exit")
				<-done
				err := d.client.ContainerKill(context.TODO(), containerName, "SIGKILL")
				if err != nil {
					d.logger.Err(err).Msg("Failed to send sigint signal")
				}
				return
			default:
				d.logger.Info().Msg("Context error killing process")
				err := d.client.ContainerKill(context.TODO(), containerName, "SIGKILL")
				if err != nil {
					d.logger.Err(err).Msg("Failed to send kill signal")
				}
				return
			}
		case <-done:
			d.logger.Info().Msg("Process exited")
			// err := d.client.ContainerKill(ctx, containerName, "SIGKILL")
			// if err != nil {
			// 	d.logger.Err(err).Msg("Failed to send sigint signal to container")
			// }
			return
		}
	}
}

func (d *DockerProvider) writeFiles(ctx context.Context, containerName string, dir string, execReq db.Jobqueue) error {
	files := map[string]string{
		"flake.nix":               execReq.Flake.String,
		execReq.ScriptPath.String: execReq.Script.String,
	}

	tarFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.tar", execReq.ID))
	d.logger.Info().Str("Path", tarFilePath).Msg("Writing files to tar")
	defer os.Remove(tarFilePath)

	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		d.logger.Err(err).Msg("Failed to create tar file")
		return err
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(content)),
			Mode: 0644,
		}); err != nil {
			d.logger.Err(err).Msg("Failed to write header")
			return err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			d.logger.Err(err).Msg("Failed to write content")
			return err
		}
	}

	err = d.client.CopyToContainer(ctx, containerName, dir, tarFile, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		d.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	return nil
}

func (d *DockerProvider) updateJob(ctx context.Context, jobID int64, message string) error {
	if _, err := d.queries.UpdateJob(ctx, db.UpdateJobParams{
		ID:   jobID,
		Logs: pgtype.Text{String: message, Valid: true},
	}); err != nil {
		return err
	}

	return nil
}
