package pool

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/namesgenerator"
	"github.com/google/uuid"
	"github.com/jackc/puddle/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func K8sSandboxConstructor(ctx context.Context) (Pod, error) {
	log.Println("Creating Sandbox -0-----------------------")
	envConfig, _ := config.GetEnvConfig()
	cli, _ := GetK8sClient()

	sandboxId := uuid.NewString()

	podLabels := map[string]string{
		"sandboxId": sandboxId,
	}

	podSpec := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    envConfig.K8S_NAMESPACE,
			GenerateName: "valkyrie-sandbox" + "-",
			Labels:       podLabels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  strings.ReplaceAll(namesgenerator.GetRandomName(10), "_", ""),
					Image: envConfig.SANDBOX_IMAGE,
					Ports: []corev1.ContainerPort{
						{
							Name:          "code-server",
							ContainerPort: 9090,
						},
						{
							Name:          "agent",
							ContainerPort: 1618,
						},
					},
				},
			},
		},
	}

	pod, err := cli.CoreV1().Pods(envConfig.K8S_NAMESPACE).Create(ctx, podSpec, metav1.CreateOptions{})
	if err != nil {
		return Pod{}, fmt.Errorf("error creating pod: %+v", err)
	}

	serviceSpec := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-service", pod.Name),
			Namespace: envConfig.K8S_NAMESPACE,
		},
		Spec: corev1.ServiceSpec{
			Selector: podLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "code-server",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromString("code-server"),
				},
				{
					Name:       "agent",
					Protocol:   corev1.ProtocolTCP,
					Port:       90,
					TargetPort: intstr.FromString("agent"),
				},
			},
		},
	}
	_, err = cli.CoreV1().Services(envConfig.K8S_NAMESPACE).Create(ctx, serviceSpec, metav1.CreateOptions{})
	if err != nil {
		return Pod{}, fmt.Errorf("error creating service: %+v", err)
	}

	return Pod{
		Name: pod.Name,
		Container: Container{
			Name: pod.Spec.Containers[0].Name,
		},
		SandboxID: sandboxId,
	}, nil
}

func K8sSandboxDestructor(pod Pod) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	envConfig, err := config.GetEnvConfig()
	if err != nil {
		log.Printf("failed to get environment config: %v", err)
	}

	cli, _ := GetK8sClient()
	if err != nil {
		log.Printf("failed to get Kubernetes client: %v", err)
	}

	serviceName := fmt.Sprintf("%s-service", pod.Name)
	deleteServiceErr := cli.CoreV1().Services(envConfig.K8S_NAMESPACE).Delete(ctx, serviceName, metav1.DeleteOptions{})
	if deleteServiceErr != nil {
		fmt.Printf("failed to delete service %s for pod %s: %v\n", serviceName, pod.Name, deleteServiceErr)
	}

	ingressName := fmt.Sprintf("%s-ingress", pod.Name)

	_, err = cli.NetworkingV1().Ingresses(envConfig.K8S_NAMESPACE).Get(context.TODO(), ingressName, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Ingress '%s' not found in namespace '%s'. Nothing to delete.\n", ingressName, envConfig.K8S_NAMESPACE)
		} else {
			fmt.Printf("Error getting Ingress '%s' in namespace '%s': %v\n", ingressName, envConfig.K8S_NAMESPACE, err)
		}
	} else {

		deletePolicy := metav1.DeletePropagationForeground
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}

		err = cli.NetworkingV1().Ingresses(envConfig.K8S_NAMESPACE).Delete(ctx, ingressName, deleteOptions)
		if err != nil {
			fmt.Printf("Error deleting Ingress '%s' in namespace '%s': %v\n", ingressName, envConfig.K8S_NAMESPACE, err)
		}

		fmt.Printf("Ingress '%s' in namespace '%s' deleted successfully.\n", ingressName, envConfig.K8S_NAMESPACE)
	}

	// Delete the Pod
	deletePodErr := cli.CoreV1().Pods(envConfig.K8S_NAMESPACE).Delete(ctx, pod.Name, metav1.DeleteOptions{})
	if deletePodErr != nil {
		log.Printf("failed to delete pod %s: %v", pod.Name, deletePodErr)
	}

	fmt.Printf("Successfully deleted pod %s and service %s\n", pod.Name, serviceName)
}

func NewK8sSandboxPool(ctx context.Context, initPoolSize int32, maxPoolSize int32) (*puddle.Pool[Pod], error) {
	pool, err := puddle.NewPool(&puddle.Config[Pod]{Constructor: K8sSandboxConstructor, Destructor: K8sSandboxDestructor, MaxSize: maxPoolSize})
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
