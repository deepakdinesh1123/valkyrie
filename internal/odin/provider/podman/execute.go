//go:build linux

package podman

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/provider/common"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/docker/docker/api/types/container"
	spec "github.com/opencontainers/runtime-spec/specs-go"
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

	// A temporary directory on the host machine to create the overlay store
	// This directory also contains a zip of the code and the flake.nix
	prepDir := filepath.Join(os.TempDir(), fmt.Sprintf("odin-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(prepDir, 0755); err != nil {
		p.logger.Err(err).Msg("Failed to create temp dir")
		return
	}

	err := common.OverlayStore(prepDir, p.envConfig.ODIN_NIX_STORE)
	if err != nil {
		p.logger.Err(err).Msg("Failed to overlay store")
		common.Cleanup(prepDir)
		return
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
		p.envConfig.ODIN_WORKER_PODMAN_IMAGE,
		false,
	)
	s.Name = containerName

	stopTimeout := uint(p.envConfig.ODIN_WORKER_TASK_TIMEOUT)
	s.StopTimeout = &stopTimeout

	stopSignal := syscall.SIGINT
	s.StopSignal = &stopSignal

	s.OCIRuntime = p.envConfig.ODIN_WORKER_RUNTIME
	s.ContainerStorageConfig.Mounts = append(s.ContainerStorageConfig.Mounts, spec.Mount{
		Destination: "/nix",
		Type:        "bind",
		Source:      filepath.Join(prepDir, "merged"),
		Options:     []string{"U", "true"},
	})

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

	var contInfo *define.InspectContainerData
	if contInfo, err = containers.Inspect(p.conn, cont.ID, nil); err != nil {
		p.logger.Err(err).Msg("Failed to inspect container")
		p.updateJob(ctx, &job, start, err.Error(), "", false)
		return
	}

	err = p.writeFiles(ctx, containerName, prepDir, job)
	if err != nil {
		p.logger.Err(err).Msg("Failed to write files")
		containers.Kill(p.conn, string(job.JobID), nil)
		p.updateJob(ctx, &job, start, err.Error(), "", false)
		return
	}

	p.logger.Info().Msg("Files written")

	done := make(chan bool, 1)
	go func() {
		execId, err := containers.ExecCreate(p.conn, cont.ID, &handlers.ExecCreateConfig{
			ExecConfig: container.ExecOptions{
				AttachStderr: true,
				AttachStdout: true,
				WorkingDir:   "/odin",
				Cmd:          []string{"nix", "run"},
			},
		})
		if err != nil {
			p.logger.Err(err).Msg("Failed to create exec")
			p.updateJob(ctx, &job, start, err.Error(), "", false)
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
				common.Cleanup(prepDir)
				common.KillContainer(contInfo.State.Pid)
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

func (p *PodmanProvider) writeFiles(ctx context.Context, containerName string, prepDir string, job db.Job) error {
	execReq, err := p.queries.GetExecRequest(ctx, job.ExecRequestID.Int32)
	if err != nil {
		return err
	}
	files := map[string]string{
		"flake.nix":  execReq.Flake,
		execReq.Path: execReq.Code,
	}

	tarFilePath, err := common.CreateTarArchive(files, prepDir)
	if err != nil {
		return err
	}

	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		p.logger.Err(err).Msg("Failed to open tar file")
		return err
	}
	defer tarFile.Close()
	defer os.Remove(tarFilePath)

	copyF, err := containers.CopyFromArchive(p.conn, containerName, "/odin", tarFile)
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
