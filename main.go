package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespace := flag.String("namespace", "default", "namespace")
	labelSelectors := flag.String("labels", "app.kubernetes.io/name=ui", "Lables to use as labelSelectors")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podsClient := clientset.CoreV1().Pods(*namespace)
outer:
	for counter := 1; counter <= 5; counter++ {
		list, err := podsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: *labelSelectors})
		if err != nil {
			panic(err)
		}

		if len(list.Items) == 0 {
			fmt.Printf("There were no pods found matching the selected criteria\n")
			os.Exit(1)
		}
		for _, d := range list.Items {
			currentStatus := d.Status.Phase
			if currentStatus == "Running" {
				fmt.Printf("Pod: %v\nStatus: %v", d.Name, currentStatus)
				break outer
			}
		}
		fmt.Printf("Waiting for targeted pod to become in runing state...\n")
		time.Sleep(5 * time.Second)
	}
}
