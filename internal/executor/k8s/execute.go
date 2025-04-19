package k8s

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
	done := make(chan bool)

	var err error

	pod, err := ke.pool.Acquire(ctx)
	if err != nil {

	}

	execReq, err := ke.queries.GetExecRequest(ctx, job.Arguments.ExecConfig.ExecReqId)

	// Write Files
	script, spec, err := execution.ConvertExecSpecToNixScript(ctx, &execReq, ke.queries)
	if err != nil {
		return
	}
	files := map[string]string{
		"exec.sh":       script,
		spec.ScriptName: execReq.Code.String,
		"input.txt":     execReq.Input.String,
	}

	go func() {
		defer func() {
			done <- true
		}()

		for {
			select {
			case <-ctx.Done():
				switch ctx.Err() {
				case context.DeadlineExceeded:
					req := ke.client.CoreV1().RESTClient().
						Post(). // Use Post method for exec
						Resource("pods").
						Name(pod.Value().Name).
						Namespace(ke.envConfig.K8S_NAMESPACE).
						SubResource("exec")

					req.VersionedParams(&corev1.PodExecOptions{
						Container: pod.Value().Container.Name,
						Command:   []string{"sh", "/home/valnix/nix_stop.sh"},
						Stdin:     false, // No stdin for this example
						Stdout:    false, // Capture stdout
						Stderr:    false, // Capture stderr
						TTY:       false, // No TTY needed for non-interactive commands
					}, scheme.ParameterCodec)

					var stdout bytes.Buffer
					streamOptions := remotecommand.StreamOptions{
						Stdout: &stdout,
					}
					exec, err := remotecommand.NewSPDYExecutor(ke.config, "POST", req.URL())
					if err != nil {
						err = fmt.Errorf("Error creating SPDY executor: %+v\n", err)
					}
					err = exec.StreamWithContext(ctx, streamOptions)
					if err != nil {

					}
					stdoutStr := stdout.String()
					_, err = ke.queries.UpdateJobResultTx(context.TODO(), db.UpdateJobResultTxParams{
						Job:      job,
						WorkerId: ke.workerId,
						ExecLogs: stdoutStr,
						Success:  false,
						Retry:    true,
					})
					ke.handleFailure(job, err)
				case context.Canceled:
					_, err = ke.queries.UpdateJobResultTx(context.TODO(), db.UpdateJobResultTxParams{
						Job:      job,
						WorkerId: ke.workerId,
						ExecLogs: "",
						Success:  false,
						Retry:    false,
					})
					ke.handleFailure(job, err)
				}
				pod.Destroy()
			default:
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
					}, scheme.ParameterCodec)

					exec, err := remotecommand.NewSPDYExecutor(ke.config, "POST", req.URL())
					if err != nil {
						err = fmt.Errorf("error creating executor: %+v", err)
						return
					}

					// Create a stream to the container
					err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
						Stdin: strings.NewReader(content),
						Tty:   false,
					})
				}

				// Execute

				req := ke.client.CoreV1().RESTClient().
					Post(). // Use Post method for exec
					Resource("pods").
					Name(pod.Value().Name).
					Namespace(ke.envConfig.K8S_NAMESPACE).
					SubResource("exec")

				req.VersionedParams(&corev1.PodExecOptions{
					Container: pod.Value().Container.Name,
					Command:   []string{"sh", "/home/valnix/nix_run.sh"},
					Stdin:     false, // No stdin for this example
					Stdout:    false, // Capture stdout
					Stderr:    false, // Capture stderr
					TTY:       false, // No TTY needed for non-interactive commands
				}, scheme.ParameterCodec)

				exec, err := remotecommand.NewSPDYExecutor(ke.config, "POST", req.URL())
				if err != nil {
					err = fmt.Errorf("Error creating SPDY executor: %+v\n", err)
					ke.handleFailure(job, err)
					return
				}

				err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{})
				if err != nil {
					ke.handleFailure(job, err)
					return
				}

				req.VersionedParams(&corev1.PodExecOptions{
					Container: pod.Value().Container.Name,
					Command:   []string{"sh", "-c", "cat /home/valnix/valkyrie/output.txt"},
					Stdin:     false, // No stdin for this example
					Stdout:    true,  // Capture stdout
					Stderr:    true,  // Capture stderr
					TTY:       false, // No TTY needed for non-interactive commands
				}, scheme.ParameterCodec)

				var stdout, stderr bytes.Buffer
				streamOptions := remotecommand.StreamOptions{
					Stdout: &stdout,
					Stderr: &stderr,
				}
				err = exec.StreamWithContext(ctx, streamOptions)
				if err != nil {
					ke.handleFailure(job, err)
					return
				}
				stdoutStr := stdout.String()

				_, err = ke.queries.UpdateJobResultTx(context.TODO(), db.UpdateJobResultTxParams{
					Job:      job,
					WorkerId: ke.workerId,
					ExecLogs: stdoutStr,
					Success:  true,
					Retry:    false,
				})
				ke.handleFailure(job, err)
			}
		}
	}()

	ke.handleFailure(job, err)
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
