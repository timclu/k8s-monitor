package main

import (
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {

	commonConfigName := "common-config"
	verLable := "common-config-ver"
	statusLable := "common-config-status"
	namespace := "tyd"

	phase := map[string]string{
		"alpha": "./config",
		"beta":  "./b_config",
		"prod":  "./p_config",
	}

	argsWithProg := os.Args
	fmt.Println("[ConfigMap Monitoring]")
	kubeconfig := flag.String("kubeconfig", phase[argsWithProg[1]], "absolute path to the kubeconfig file")
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	commonConfig, err := clientset.CoreV1().ConfigMaps(namespace).Get(commonConfigName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	configName := commonConfig.GetName()
	configVer := commonConfig.Labels[verLable]
	fmt.Printf("ConfigMap:\n\t%s: %s\n", configName, configVer)

	result, err := clientset.ExtensionsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for i := range result.Items {
		status := result.Items[i].Labels[statusLable]
		if status == "enabled" {
			name := result.Items[i].GetName()
			ver := result.Items[i].Labels[verLable]
			var podName []string

			fmt.Println("************************************")
			fmt.Printf("%s:\n", name)

			podResult, err := clientset.Core().Pods(namespace).List(
				metav1.ListOptions{
					LabelSelector: "name=" + name,
				})
			if err != nil {
				panic(err.Error())
			}

			podName = make([]string, len(podResult.Items))

			for j := range podResult.Items {
				podName[j] = podResult.Items[j].GetName()
				fmt.Printf("\tpod[%d]:%s\n", j+1, podName[j])
			}

			fmt.Printf("\tcommon-config-ver:%s\n", ver)
			if ver == configVer {
				fmt.Printf("\t(Version Match!)\n")
			} else {
				fmt.Printf("\t(Version not match!)\n")
				fmt.Printf("\t#start to restart the pod#\n")
			}
		}
	}
	//options := metav1.DeleteOptions{}
}
