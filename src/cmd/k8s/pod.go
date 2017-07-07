package k8s

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//GetPod - get pod list
func GetPod(clientset *kubernetes.Clientset, namespace string, options metav1.ListOptions) *v1.PodList {
	podResult, err := clientset.Core().Pods(namespace).List(options)
	if err != nil {
		panic(err.Error())
	}

	return podResult
}
