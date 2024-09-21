//go:build linux

package podman

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/docker/docker/api/types/container"
)

func (p *PodmanProvider) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job db.Job) {
	start := time.Now()

	var timeout int
	if job.TimeOut.Int32 > 0 { // By default, timeout is set to -1
		timeout = int(job.TimeOut.Int32)
	} else if job.TimeOut.Int32 == 0 {
		timeout = 0
	} else {
		timeout = p.envConfig.ODIN_WORKER_TASK_TIMEOUT
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
	s := specgen.NewSpecGenerator(
		"docker.io/deepakdinesh/nix:alpine_amd64",
		false,
	)
	s.Name = containerName

	stopTimeout := uint(p.envConfig.ODIN_WORKER_TASK_TIMEOUT)
	s.StopTimeout = &stopTimeout

	stopSignal := syscall.SIGINT
	s.StopSignal = &stopSignal

	s.OCIRuntime = p.envConfig.ODIN_WORKER_RUNTIME

	containerRemove := false
	s.Remove = &containerRemove
	p.logger.Info().Msg("Creating container spec")
	cont, err := containers.CreateWithSpec(p.conn, s, nil)
	if err != nil {
		p.logger.Err(err)
		p.updateJob(ctx, &job, start, err.Error(), "", false)
	}
	p.logger.Info().Str("Container ID", cont.ID).Msg("Container created")

	err = containers.Start(p.conn, cont.ID, nil)
	if err != nil {
		p.logger.Err(err).Msg("Could not start container")
	}
	p.logger.Info().Msg("Container started")

	if _, err := containers.Inspect(p.conn, cont.ID, nil); err != nil {
		p.logger.Err(err).Msg("Failed to inspect container")
		p.updateJob(ctx, &job, start, err.Error(), "", false)
		return
	}

	err = p.writeFiles(ctx, containerName, job)
	if err != nil {
		p.logger.Err(err).Msg("Failed to write files")
		containers.Kill(p.conn, string(job.JobID), nil)
		p.updateJob(ctx, &job, start,   err.Error(), "", false)
		return
	}

	p.logger.Info().Msg("Files written")

	done := make(chan bool, 1)
	go func() {
		execId, err := containers.ExecCreate(p.conn, cont.ID, &handlers.ExecCreateConfig{
			ExecConfig: container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				WorkingDir:   "/home/valnix/odin",
				Cmd:          []string{"nix", "run"},
			},
		})
		if err != nil {
			p.logger.Err(err).Msg("Failed to create exec")
			p.updateJob(ctx, &job,start, err.Error(), "", false)
			done <- true
			return
		}
		p.logger.Info().Msg("Exec created")

		r, w, err := os.Pipe()
		if err != nil {
			p.logger.Err(err).Msg("Could not create pipe")
		}
		defer r.Close()
		execOpts := new(containers.ExecStartAndAttachOptions).WithOutputStream(w).WithAttachOutput(true).WithErrorStream(w)
		err = containers.ExecStartAndAttach(p.conn, execId, execOpts)
		if err != nil {
			p.logger.Err(err).Msg("Failed to start exec")
			p.updateJob(ctx, &job, start, err.Error(), "", false)
			done <- true
			return
		}
		outputC := make(chan string)

		// copying output in a separate goroutine, so that printing doesn't remain blocked forever
		go func() {
			var output bytes.Buffer
			_, _ = io.Copy(&output, r)
			outputC <- output.String()
		}()
		w.Close()
		p.updateJob(ctx, &job, start, "Completed", <-outputC, true)
		done <- true
	}()

	for {
		select {
		case <-tctx.Done():
			switch tctx.Err() {
			case context.DeadlineExceeded:
				p.logger.Info().Msg("Context deadline exceeded wating for process to exit")
				err := containers.Kill(p.conn, cont.ID, new(containers.KillOptions).WithSignal("SIGINT"))
				if err != nil {
					p.logger.Err(err).Msg("Failed to send sigint signal")
				}
				return
			}
		case <-ctx.Done():
			switch ctx.Err() {
			case context.Canceled:
				p.logger.Info().Msg("Time out killing process")
				<-done
				err := containers.Kill(p.conn, cont.ID, new(containers.KillOptions).WithSignal("SIGINT"))
				if err != nil {
					p.logger.Err(err).Msg("Failed to send sigint signal")
				}
				return
			default:
				p.logger.Info().Msg("Context error killing process")
				<-done
				err := containers.Kill(p.conn, cont.ID, new(containers.KillOptions).WithSignal("SIGKILL"))
				if err != nil {
					p.logger.Err(err).Msg("Failed to send sigkill signal")
				}
				return
			}
		case <-done:
			p.logger.Info().Msg("Process exited")
			err := containers.Kill(p.conn, cont.ID, new(containers.KillOptions).WithSignal("SIGKILL"))
			if err != nil {
				p.logger.Err(err).Msg("Failed to send sigkill signal")
			}
			return
		}
	}

}

func (p *PodmanProvider) updateJob(ctx context.Context, job *db.Job, startTime time.Time, message string, logs string, success bool) error {
	p.logger.Info().Bool("Success", success).Str("Message", message).Str("Logs", logs).Msg("Updating job result")
	retry := true
	if job.Retries.Int32+1 >= job.MaxRetries.Int32 || success {
		retry = false
	}
	if _, err := p.queries.UpdateJobResultTx(ctx, db.UpdateJobResultTxParams{
		StartTime: startTime,
		Job:       *job,
		Message:   message,
		Success:   success,
		Retry:     retry,
		WorkerId:  p.workerId,
	}); err != nil {
		return err
	}
	return nil
}

func (p *PodmanProvider) writeFiles(ctx context.Context, containerName string, job db.Job) error {
	execReq, err := p.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
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
		p.logger.Err(err).Msg("Failed to open tar file")
		return err
	}
	defer tarFile.Close()

	copyF, err := containers.CopyFromArchive(p.conn, containerName, "/home/valnix/odin/", tarFile)
	if err != nil {
		p.logger.Err(err).Msg("Failed to copy files to container")
		return err
	}
	err = copyF()
	if err != nil {
		p.logger.Err(err).Msg("Failed to copy files to container")
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
