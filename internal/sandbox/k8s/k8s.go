package k8s

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/deepakdinesh1123/valkyrie/internal/concurrency"
	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/db/jsonschema"
	"github.com/deepakdinesh1123/valkyrie/internal/pool"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/puddle/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type K8SH struct {
	client      *kubernetes.Clientset
	config      *rest.Config
	queries     db.Store
	envConfig   *config.EnvConfig
	workerId    int32
	logger      *zerolog.Logger
	tp          trace.TracerProvider
	mp          metric.MeterProvider
	sandboxPool *puddle.Pool[pool.Pod]

	pods []*puddle.Resource[pool.Pod]
	mu   sync.Mutex
}

func NewK8SandboxHandler(ctx context.Context, queries db.Store, workerId int32, tp trace.TracerProvider, mp metric.MeterProvider, envConfig *config.EnvConfig, logger *zerolog.Logger) (*K8SH, error) {
	client, config := pool.GetK8sClient()
	if client == nil {
		return nil, fmt.Errorf("could not get docker client")
	}
	return &K8SH{
		client:    client,
		config:    config,
		queries:   queries,
		envConfig: envConfig,
		workerId:  workerId,
		logger:    logger,
		tp:        tp,
		mp:        mp,
	}, nil
}

func (k *K8SH) StartSandboxPool(ctx context.Context, envConfig *config.EnvConfig) error {
	sandboxPool, err := pool.NewK8sSandboxPool(ctx, int32(envConfig.HOT_CONTAINER), int32(envConfig.WORKER_CONCURRENCY))
	if err != nil {
		return fmt.Errorf("error creating sandbox pool: %s", err)
	}
	k.sandboxPool = sandboxPool
	return nil
}

func (k *K8SH) Create(ctx context.Context, wg *concurrency.SafeWaitGroup, sandBoxJob db.FetchSandboxJobTxResult) {
	sandbox := sandBoxJob.Sandbox
	handleError := func(err error, message string) {
		if err != nil {
			k.logger.Err(err).Msg(message)
			// Mark sandbox as failed in database
			updateErr := k.queries.UpdateSandboxState(ctx, db.UpdateSandboxStateParams{
				SandboxID:    sandbox.SandboxID,
				CurrentState: "failed",
				Details:      jsonschema.SandboxDetails{Error: err.Error()},
			})
			if updateErr != nil {
				k.logger.Err(updateErr).Msg("failed to mark sandbox as failed")
			}
			return
		}
	}
	pod, err := k.sandboxPool.Acquire(ctx)
	if err != nil {
		handleError(err, "could not acquire pod")
		return
	}

	_, err = k.createIngress(ctx, pod)
	if err != nil {
		handleError(err, "could not create ingress")
		return
	}

	if sandbox.Config.Flake != "" {
		// Write the flake to the sandbox
		flakeCmd := []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("rm -f /home/valnix/work/flake.nix && echo \"%s\" >> /home/valnix/work/flake.nix && nix profile remove work", sandbox.Config.Flake),
		}

		stdout, stderr, err := k.executeCmdInPod(ctx, pod.Value().Name, pod.Value().Container.Name, flakeCmd)
		if err != nil {
			k.logger.Err(err).Msgf("error adding flake in sandbox pod %s: %v\nstdout: %s\nstderr: %s",
				pod.Value().Name, err, stdout, stderr)
			handleError(err, "failed to add flake to pod")
			return
		}

		// Create a command to install the flake
		installCmd := []string{
			"nix",
			"profile",
			"install",
			".",
		}

		stdout, stderr, err = k.executeCmdInPod(ctx, pod.Value().Name, pod.Value().Container.Name, installCmd, "/home/valnix/work")
		if err != nil {
			k.logger.Err(err).Msgf("error installing flake in sandbox pod %s: %v\nstdout: %s\nstderr: %s",
				pod.Value().Name, err, stdout, stderr)
			handleError(err, "failed to install flake in pod")
			return
		}
	}

	containerURL := fmt.Sprintf("http:/%s-cs.%s", pod.Value().Name, k.envConfig.SANDBOX_HOSTNAME)
	sandboxAgentUrl := fmt.Sprintf("ws:/%s-ag.%s/sandbox", pod.Value().Name, k.envConfig.SANDBOX_HOSTNAME)

	k.logger.Info().Msg(containerURL)
	err = k.queries.MarkSandboxRunning(ctx, db.MarkSandboxRunningParams{
		SandboxID:       sandbox.SandboxID,
		SandboxUrl:      pgtype.Text{String: containerURL, Valid: true},
		SandboxAgentUrl: pgtype.Text{String: sandboxAgentUrl, Valid: true},
	})
	if err != nil {
		handleError(err, "error marking the container as running")
		return
	}

	k.mu.Lock()
	k.pods = append(k.pods, pod)
	k.mu.Unlock()

	err = k.queries.UpdateJobCompleted(ctx, sandBoxJob.Job.JobID)
	if err != nil {
		k.logger.Err(err).Msgf("error changing job state: %d", sandBoxJob.Job.JobID)
	}

}

func (k *K8SH) createIngress(ctx context.Context, pod *puddle.Resource[pool.Pod]) (*networkingv1.Ingress, error) {
	pathTypePrefix := networkingv1.PathTypePrefix
	ingresSpec := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-ingress", pod.Value().Name),
			Namespace: k.envConfig.K8S_NAMESPACE,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: fmt.Sprintf("%s-cs.%s", pod.Value().Name, k.envConfig.SANDBOX_HOSTNAME),
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: fmt.Sprintf("%s-service", pod.Value().Name),
											Port: networkingv1.ServiceBackendPort{
												Name: "code-server",
											},
										},
									},
									PathType: &pathTypePrefix,
								},
							},
						},
					},
				},
				{
					Host: fmt.Sprintf("%s-ag.%s", pod.Value().Name, k.envConfig.SANDBOX_HOSTNAME),
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: fmt.Sprintf("%s-service", pod.Value().Name),
											Port: networkingv1.ServiceBackendPort{
												Name: "agent",
											},
										},
									},
									PathType: &pathTypePrefix,
								},
							},
						},
					},
				},
			},
		},
	}
	ing, err := k.client.NetworkingV1().Ingresses(k.envConfig.K8S_NAMESPACE).Create(ctx, ingresSpec, metav1.CreateOptions{})
	return ing, err
}

func (k *K8SH) executeCmdInPod(ctx context.Context, podName, containerName string, command []string, workingDir ...string) (string, string, error) {
	req := k.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(k.envConfig.K8S_NAMESPACE).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	if len(workingDir) > 0 && workingDir[0] != "" {
		command = append([]string{"cd", workingDir[0], "&&"}, command...)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	var stdout, stderr bytes.Buffer
	exec, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error creating SPDY executor: %v", err)
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	return stdout.String(), stderr.String(), err
}

func (d *K8SH) Cleanup(ctx context.Context) error {
	for _, cont := range d.pods {
		cont.Destroy()
	}
	d.sandboxPool.Close()
	return nil
}
