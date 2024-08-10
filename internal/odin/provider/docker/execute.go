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

func (d *DockerProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, execReq db.Job) {
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
		err := d.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	err = d.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		d.logger.Err(err).Msg("Failed to start container")
		err := d.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}

	if _, err := d.client.ContainerInspect(ctx, resp.ID); err != nil {
		d.logger.Err(err).Msg("Failed to inspect container")
		err := d.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}

	err = d.writeFiles(ctx, containerName, execReq)
	if err != nil {
		d.logger.Err(err).Msg("Failed to write files")
		d.client.ContainerKill(ctx, containerName, "SIGKILL")
		err := d.updateJob(ctx, &execReq, err.Error())
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
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
			err := d.updateJob(ctx, &execReq, err.Error())
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		hijResp, err := d.client.ContainerExecAttach(ctx, resp.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msg("Failed to attach exec")
			err := d.updateJob(ctx, &execReq, err.Error())
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		out, err := io.ReadAll(hijResp.Reader)
		if err != nil {
			d.logger.Err(err).Msg("Failed to read output")
			err := d.updateJob(ctx, &execReq, err.Error())
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		err = d.updateJob(ctx, &execReq, stripCtlAndExtFromUTF8(string(out)))
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
				d.logger.Info().Msg("Context canceled, waiting for processes to finish")
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

func (d *DockerProvider) writeFiles(ctx context.Context, containerName string, execReq db.Job) error {
	files := map[string]string{
		"flake.nix":        execReq.Flake,
		execReq.ScriptPath: execReq.Script,
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

// updateJob updates the job status to completed and inserts a new job run.
//
// Parameters:
// - ctx: the context for the update operation.
// - execReq: the job execution request.
// - message: the message to be logged.
// Returns:
// - error: an error if the update operation fails.
func (d *DockerProvider) updateJob(ctx context.Context, execReq *db.Job, message string) error {
	if err := d.queries.UpdateJob(ctx, execReq.ID); err != nil {
		return err
	}
	if _, err := d.queries.InsertJobRun(ctx, db.InsertJobRunParams{
		JobID:      execReq.ID,
		WorkerID:   execReq.WorkerID.Int32,
		StartedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		FinishedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Script:     message,
		Flake:      execReq.Flake,
		Args:       execReq.Args,
		Logs:       pgtype.Text{String: message, Valid: true},
	}); err != nil {
		return err
	}
	return nil
}

// stripCtlAndExtFromUTF8 removes control characters and non-ASCII characters from a UTF8 string.
//
// Parameters:
// - str: the input string to be processed.
// Returns:
// - string: the processed string with control characters and non-ASCII characters removed.
func stripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 || r == 10 {
			return r
		}
		return -1
	}, str)
}
