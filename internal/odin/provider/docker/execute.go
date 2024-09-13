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
)

func (d *DockerProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job) {
	start := time.Now()
	var timeout int
	if job.TimeOut.Int32 > 0 { // By default, timeout is set to -1
		timeout = int(job.TimeOut.Int32)
	} else if job.TimeOut.Int32 == 0 {
		timeout = 0
	} else {
		timeout = d.envConfig.ODIN_WORKER_TASK_TIMEOUT
	}

	var tctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		tctx, cancel = context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	} else {
		tctx, cancel = context.WithCancel(context.TODO())
	}
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
		err := d.updateJob(ctx, &job, start, err.Error(), false)
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}
	err = d.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		d.logger.Err(err).Msg("Failed to start container")
		err := d.updateJob(ctx, &job, start, err.Error(), false)
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}

	if _, err := d.client.ContainerInspect(ctx, resp.ID); err != nil {
		d.logger.Err(err).Msg("Failed to inspect container")
		err := d.updateJob(ctx, &job, start, err.Error(), false)
		if err != nil {
			d.logger.Err(err).Msg("Failed to update job")
		}
		return
	}

	err = d.writeFiles(ctx, containerName, job)
	if err != nil {
		d.logger.Err(err).Msg("Failed to write files")
		d.client.ContainerKill(ctx, containerName, "SIGKILL")
		err := d.updateJob(ctx, &job, start, err.Error(), false)
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
			err := d.updateJob(ctx, &job, start, err.Error(), false)
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		hijResp, err := d.client.ContainerExecAttach(ctx, resp.ID, container.ExecAttachOptions{})
		if err != nil {
			d.logger.Err(err).Msg("Failed to attach exec")
			err := d.updateJob(ctx, &job, start, err.Error(), false)
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		out, err := io.ReadAll(hijResp.Reader)
		if err != nil {
			d.logger.Err(err).Msg("Failed to read output")
			err := d.updateJob(ctx, &job, start, err.Error(), false)
			if err != nil {
				d.logger.Err(err).Msg("Failed to update job")
			}
			done <- true
			return
		}
		err = d.updateJob(ctx, &job, start, stripCtlAndExtFromUTF8(string(out)), true)
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

func (d *DockerProvider) writeFiles(ctx context.Context, containerName string, job db.Job) error {
	execReq, err := d.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		return err
	}
	files := map[string]string{
		"flake.nix":  execReq.Flake,
		execReq.Path: execReq.Code,
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

func (d *DockerProvider) updateJob(ctx context.Context, job *db.Job, startTime time.Time, message string, success bool) error {
	retry := true
	if job.Retries.Int32+1 >= job.MaxRetries.Int32 || success {
		retry = false
	}
	if _, err := d.queries.UpdateJobResultTx(ctx, db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       *job,
		Message:   message,
		Success:   success,
		Retry:     retry,
		WorkerId:  d.workerId,
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
