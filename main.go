package main

import (
	"flag"
	"fmt"
	"log"
	"pedro-gutierrez/voting/pkg/clienthelper"
	"pedro-gutierrez/voting/pkg/podrequirements"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	var kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file (if not provided, will run in 'incluster' mode)")
	var namespace = flag.String("namespace", "", "the namespace where to get pods from")

	flag.Parse()

	// get a client to the kubernetes api server
	clientset, err := clienthelper.NewClientset(*kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// get the list of pods for the given namespace
	pods, err := clientset.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	// obtained aggregated cpu and memory for each pod
	// then convert it into a JSON string
	json, err := podrequirements.GetPodsCPUAndMemoryRequirements(pods).ToJSON()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s", json)
}
