package main

import (
	k8s "cmd/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt"
	"os"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	commonConfigName := "common-config"
	verLable := "common-config-ver"
	statusLable := "common-config-status"
	namespace := "tyd"
	defaultPhase := []string{"default", "alpha"}

	//phase selection
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Printf("Set phase to default stage alpha\n")
		argsWithProg = defaultPhase
	} 

	if(argsWithProg[1]!="alpha" && argsWithProg[1]!="beta" && argsWithProg[1]!="prod"){
		fmt.Printf("phase not correct!!!\n")
		argsWithProg = defaultPhase
	}

	//get k8s client
	fmt.Println("[ConfigMap Monitoring]")
	clientset := k8s.SetClientset(argsWithProg[1])

	//get configmap by name
	commonConfig := k8s.GetConfigmap(clientset, namespace, commonConfigName, verLable)
	configName := commonConfig.GetName()
	configVer := commonConfig.Labels[verLable]
	fmt.Printf("ConfigMap:\n\t%s: %s\n", configName, configVer)

	//get all deployments list
	deploymentList := k8s.GetDeploymentList(clientset, namespace)
	deployInfo := k8s.ConfigVerCompared(deploymentList, configVer, statusLable, verLable)

	//List the config info
	for i := range deployInfo {
		name := deployInfo[i].Name
		ver := deployInfo[i].Ver
		fmt.Println("************************************")
		fmt.Printf("%s:\n", name)
		fmt.Printf("\tcommon-config-ver:%s", ver)
		if ver == configVer {
			fmt.Printf(" (Version Match!)\n")
		} else {
			fmt.Printf(" (Version not match)\n")
			fmt.Printf("\tPod List:\n")
			podResult := k8s.GetPod(
				clientset,
				namespace,
				metav1.ListOptions{
					LabelSelector: "name=" + name,
				})
			podName := make([]string, len(podResult.Items))
			for j := range podResult.Items {
				podName[j] = podResult.Items[j].GetName()
				fmt.Printf("\t\tpod[%d]:%s\n", j+1, podName[j])
			}
		}
	}

}
