package config

import (
	"github.com/mclenhard/containiq/pkg/setup"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type Framework struct {
	Config *setup.Config
	Client *kubernetes.Clientset
}

func NewFramework() (*Framework , error){
	home := os.Getenv("HOME")
	path := home + "/.kube/config"
	Framework := Framework{}
	c := setup.Setup()

	config, err := clientcmd.BuildConfigFromFlags("",path)
	if err != nil {
		return nil, err
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		return nil, err
	}


	Framework.Config = c
	Framework.Client = clientset
	return &Framework, nil
}

