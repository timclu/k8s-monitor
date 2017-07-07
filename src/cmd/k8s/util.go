package k8s

import (
	"model"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//SetClientset - init the k8s client
func SetClientset(phaseConf string) *kubernetes.Clientset {
	phase := map[string]string{
		"alpha": "./config",
		"beta":  "./b_config",
		"prod":  "./p_config",
	}
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", phase[phaseConf])
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

//ConfigVerCompared - used to check common parameter status
func ConfigVerCompared(deploymentList *v1beta1.DeploymentList, configVer string, statusLabel string, verLabel string) []model.DeploymentInfo {
	var infOutput []model.DeploymentInfo
	for i := range deploymentList.Items {
		status := deploymentList.Items[i].Labels[statusLabel]
		if status == "enabled" {
			name := deploymentList.Items[i].GetName()
			ver := deploymentList.Items[i].Labels[verLabel]
			infOutput = append(infOutput, model.DeploymentInfo{name, ver})
		}
	}
	return infOutput
}
