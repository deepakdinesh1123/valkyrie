package k8s

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/deepakdinesh1123/valkyrie/internal/services/execution"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sExecutor struct {
	queries   db.Store
	envConfig *config.EnvConfig
	workerId  int32
	logger    *zerolog.Logger
	tp        trace.TracerProvider
	mp        metric.MeterProvider
	pool      *puddle.Pool[pool.Pod]
	client    *kubernetes.Clientset
	config    *rest.Config
}

func NewK8sExecutor(ctx context.Context, env *config.EnvConfig, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, logger *zerolog.Logger) (*K8sExecutor, error) {
	client, config := pool.GetK8sClient()
	k8sPool, err := pool.NewK8sExecutionPool(ctx, int32(env.HOT_CONTAINER), env.WORKER_CONCURRENCY)
	if err != nil {
		return nil, fmt.Errorf("error creating executor +%v", err)
	}

	return &K8sExecutor{
		envConfig: env,
		queries:   queries,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
		client:    client,
		pool:      k8sPool,
		config:    config,
	}, nil
}

func (ke *K8sExecutor) Execute(ctx context.Context, wg *concurrency.SafeWaitGroup, job *db.Job, logger zerolog.Logger) {
	defer wg.Done()

	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	pod, err := ke.pool.Acquire(execCtx)
	if err != nil {
		ke.handleFailure(job, fmt.Errorf("could not acquire pod: %w", err))
		return
	}
	defer pod.Destroy()
	go ke.pool.CreateResource(ctx)

	execReq, err := ke.queries.GetExecRequest(execCtx, job.Arguments.ExecConfig.ExecReqId)
	if err != nil {
		ke.handleFailure(job, fmt.Errorf("error getting exec request: %w", err))
		return
	}

	script, spec, err := execution.ConvertExecSpecToScript(execCtx, &execReq, ke.queries)
	if err != nil {
		ke.handleFailure(job, fmt.Errorf("error converting exec spec to nix script: %w", err))
		return
	}

	files := map[string]string{
		"exec.sh":       script,
		spec.ScriptName: execReq.Code.String,
		"input.txt":     execReq.Input.String,
	}

	// Write files to the pod
	err = ke.writeFilesToPod(execCtx, pod, files)
	if err != nil {
		ke.handleExecutionResult(job, err, err.Error(), true)
		return
	}

	// Check if context was cancelled during file writing
	if execCtx.Err() != nil {
		ke.handleExecutionResult(job, execCtx.Err(), execCtx.Err().Error(), execCtx.Err() == context.DeadlineExceeded)
		return
	}

	// Run script execution and output collection concurrently
	var runErr, outputErr error

	// Create a channel to signal execution completion
	execDone := make(chan struct{})

	// Execute the script
	go func() {

		_, _, err := ke.executeCommand(execCtx, pod, []string{"sh", "/home/valnix/nix_run.sh"})
		if err != nil {
			runErr = fmt.Errorf("execute script failed: %w", err)
			cancel() // Cancel the context to terminate other operations
			return
		}
		ke.logger.Info().Msg("Command ran")
		// Signal that execution is complete
		close(execDone)
	}()

	// Wait for either completion or context cancellation
	select {
	case <-execDone:
		// Script completed successfully, retrieve output
		if runErr == nil {
			outputStdout, _, err := ke.executeCommand(execCtx, pod, []string{"sh", "-c", "cat /home/valnix/valkyrie/output.txt"})
			if err != nil {
				outputErr = fmt.Errorf("get output failed: %w", err)
				ke.handleExecutionResult(job, err, outputErr.Error(), true)
			} else {
				ke.handleExecutionResult(job, nil, outputStdout, true)
			}
		} else {
			ke.handleExecutionResult(job, runErr, runErr.Error(), true)
		}

	case <-ctx.Done():

		// Run cleanup in a separate context
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cleanupCancel()

		_, _, err := ke.executeCommand(cleanupCtx, pod, []string{"sh", "/home/valnix/nix_stop.sh"})
		if err != nil {
			ke.logger.Err(err).Msg("cleanup: could not execute command")
		} else {
			ke.logger.Info().Msg("Cleanup script executed successfully.")
		}

		ke.logger.Info().Msg("Command execution failed adding result")
		ke.handleExecutionResult(job, ctx.Err(), fmt.Sprintf("Job cancelled. Context error: %v\n", ctx.Err()), ctx.Err() == context.DeadlineExceeded)
	}
}

// writeFilesToPod writes files to the pod and returns any error
func (ke *K8sExecutor) writeFilesToPod(ctx context.Context, pod *puddle.Resource[pool.Pod], files map[string]string) error {
	for fileName, content := range files {
		req := ke.client.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(pod.Value().Name).
			Namespace(ke.envConfig.K8S_NAMESPACE).
			SubResource("exec").
			Param("container", pod.Value().Container.Name)

		req.VersionedParams(&corev1.PodExecOptions{
			Container: pod.Value().Container.Name,
			Command:   []string{"bash", "-c", "cat > " + filepath.Join("/home/valnix/valkyrie", fileName)},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

		var writeStdout, writeStderr bytes.Buffer
		exec, err := remotecommand.NewSPDYExecutor(ke.config, "POST", req.URL())
		if err != nil {
			execErr := fmt.Errorf("write file %s: could not create k8s executor: %w", fileName, err)
			return execErr
		}

		err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  strings.NewReader(content),
			Stdout: &writeStdout,
			Stderr: &writeStderr,
			Tty:    false,
		})
		if err != nil {
			execErr := fmt.Errorf("write file %s: stream command failed: %w", fileName, err)
			return execErr
		}
	}
	return nil
}

// executeCommand runs a command in the pod and returns stdout, stderr and any error
func (ke *K8sExecutor) executeCommand(ctx context.Context, pod *puddle.Resource[pool.Pod], command []string) (string, string, error) {
	req := ke.client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod.Value().Name).
		Namespace(ke.envConfig.K8S_NAMESPACE).
		SubResource("exec")

	var stdout, stderr bytes.Buffer
	req.VersionedParams(&corev1.PodExecOptions{
		Container: pod.Value().Container.Name,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(ke.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("could not create k8s executor: %w", err)
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	ke.logger.Info().Msg("Execution completed")
	ke.logger.Info().Str("stdout", stdout.String()).Str("stderr", stderr.String()).Msg("output")

	return stdout.String(), stderr.String(), err
}

// handleExecutionResult updates the job result in the database
func (ke *K8sExecutor) handleExecutionResult(job *db.Job, execErr error, execLogs string, retry bool) {
	success := execErr == nil

	if !success {
		ke.handleFailure(job, execErr)
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, dbErr := ke.queries.UpdateJobResultTx(dbCtx, db.UpdateJobResultTxParams{
		Job:      job,
		WorkerId: ke.workerId,
		ExecLogs: execLogs,
		Success:  success,
		Retry:    retry && !success, // Only retry if not successful and retry flag is set
	})

	if dbErr != nil {
		ke.logger.Err(dbErr).Msg("Failed to update job result in database")
	}
}

func (ke *K8sExecutor) handleFailure(job *db.Job, err error) {
	if err != nil {
		ke.queries.UpdateJobResultTx(context.TODO(), db.UpdateJobResultTxParams{
			Job:      job,
			WorkerId: ke.workerId,
			ExecLogs: err.Error(),
			Success:  false,
			Retry:    true,
		})
	}
}

func (ke *K8sExecutor) Cleanup() {

}
