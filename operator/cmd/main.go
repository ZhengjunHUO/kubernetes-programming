package main

import (
	"flag"
	"log"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"github.com/ZhengjunHUO/kubernetes-programming/operator/pkg/controller"
	hzjcs "github.com/ZhengjunHUO/k8s-custom-controller/pkg/client/clientset/versioned"
)

var (
	kubeconfigPath string
	namespace string
)

func main() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to kubeconfig")
	flag.StringVar(&namespace, "namespace", "default", "k8s namespace to deploy operator")
	flag.Parse()

	// Load config from kubeconfig file given in command line
	// if absent will try to load from serviceaccount's token (incluster)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Get clientset for built-in resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	hzjclientset, err := hzjcs.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Start controller
	controller.Start(clientset, hzjclientset, namespace)
}
