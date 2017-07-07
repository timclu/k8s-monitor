package k8s

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//get the deployment list
func GetDeploymentList(clientset *kubernetes.Clientset, namespace string) *v1beta1.DeploymentList {
	deploymentList, err := clientset.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return deploymentList
}
