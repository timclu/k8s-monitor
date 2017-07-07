package k8s

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//get the configmap context
func GetConfigmap(clientset *kubernetes.Clientset, namepsace string, name string, verLable string) *v1.ConfigMap {
	commonConfig, err := clientset.CoreV1().ConfigMaps(namepsace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return commonConfig
}
