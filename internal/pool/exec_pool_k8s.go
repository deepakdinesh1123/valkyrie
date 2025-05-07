package pool

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/jackc/puddle/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func K8sExecConstructor(ctx context.Context) (Pod, error) {
	envConfig, _ := config.GetEnvConfig()
	cli, _ := GetK8sClient()

	podSpec := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    envConfig.K8S_NAMESPACE,
			GenerateName: "valkyrie-exec" + "-",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  strings.ReplaceAll(namesgenerator.GetRandomName(10), "_", ""),
					Image: envConfig.EXECUTION_IMAGE,
				},
			},
		},
	}

	pod, err := cli.CoreV1().Pods(envConfig.K8S_NAMESPACE).Create(ctx, podSpec, metav1.CreateOptions{})
	if err != nil {
		return Pod{}, fmt.Errorf("error creating pod: %+v", err)
	}

	return Pod{
		Name: pod.Name,
		Container: Container{
			Name: pod.Spec.Containers[0].Name,
		},
	}, nil
}

func K8sExecDestructor(pod Pod) {
	envConfig, _ := config.GetEnvConfig()

	cli, _ := GetK8sClient()

	err := cli.CoreV1().Pods(envConfig.K8S_NAMESPACE).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("error deleting pod: %+v", err)
	}
}

func NewK8sExecutionPool(ctx context.Context, initPoolSize int32, maxPoolSize int32) (*puddle.Pool[Pod], error) {
	pool, err := puddle.NewPool(&puddle.Config[Pod]{Constructor: K8sExecConstructor, Destructor: K8sExecDestructor, MaxSize: maxPoolSize})
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(initPoolSize); i += 1 {
		err := pool.CreateResource(ctx)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
