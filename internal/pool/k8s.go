package pool

import (
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var getK8sClientOnce sync.Once
var k8sClient *kubernetes.Clientset
var restConfig *rest.Config

func GetK8sClient() (*kubernetes.Clientset, *rest.Config) {
	getK8sClientOnce.Do(
		func() {
			config, err := rest.InClusterConfig()
			if err != nil {
				log.Println("Error creating in-cluster config:", err)
				return
			}
			restConfig = config
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				log.Println("Error creating Kubernetes client:", err)
				return
			}
			k8sClient = clientset
		},
	)
	return k8sClient, restConfig
}
