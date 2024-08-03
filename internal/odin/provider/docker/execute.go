package docker

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/docker/docker/api/types/container"
	"github.com/jackc/pgx/v5/pgtype"
)

func (d *DockerProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Jobqueue) {
	tctx, cancel := context.WithTimeout(ctx, time.Duration(d.envConfig.ODIN_WORKER_TASK_TIMEOUT)*time.Second)
	defer cancel()
	defer wg.Done()
	containerName := namesgenerator.GetRandomName(0)
	resp, err := d.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:       "alpinix",
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

	err = d.writeFiles(ctx, containerName, execReq)
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
				WorkingDir:   "/home/valnix/odin",
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
		out, err := io.ReadAll(hijResp.Reader)
		if err != nil {
			d.logger.Err(err).Msg("Failed to read output")
			d.updateJob(ctx, execReq.ID, err.Error())
			done <- true
			return
		}
		err = d.updateJob(ctx, execReq.ID, stripCtlAndExtFromUTF8(string(out)))
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		done <- true
		defer hijResp.Close()
	}()

	for {
		select {
		case <-tctx.Done():
			switch tctx.Err() {
			case context.DeadlineExceeded:
				d.logger.Info().Msg("Context deadline exceeded wating for process to exit")
				err := d.client.ContainerKill(context.TODO(), containerName, "SIGKILL")
				if err != nil {
					d.logger.Err(err).Msg("Failed to send sigint signal")
				}
				return
			}
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				d.logger.Info().Msg("Time out killing process")
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
			err := d.client.ContainerKill(ctx, containerName, "SIGKILL")
			if err != nil {
				d.logger.Err(err).Msg("Failed to send sigint signal to container")
			}
			return
		}
	}
}

func (d *DockerProvider) writeFiles(ctx context.Context, containerName string, execReq db.Jobqueue) error {
	files := map[string]string{
		"flake.nix":               execReq.Flake.String,
		execReq.ScriptPath.String: execReq.Script.String,
	}

	tarFilePath, err := createTarArchive(files)
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

	err = d.client.CopyToContainer(ctx, containerName, "/home/valnix/odin/", tarFile, container.CopyToContainerOptions{})
	if err != nil {
		d.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	return nil
}

func createTarArchive(files map[string]string) (string, error) {
	tarFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.tar", time.Now().UnixNano()))
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return "", err
	}
	defer tarFile.Close()

	tw := tar.NewWriter(tarFile)
	defer tw.Close()

	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name: name,
			Size: int64(len(content)),
			Mode: 0744,
		}); err != nil {
			return "", err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return "", err
		}
	}
	return tarFilePath, nil
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

func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 || r == 10 {
			return r
		}
		return -1
	}, str)
}
