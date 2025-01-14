package ck8s

import (
	"context"

	"github.com/daniarmas/notes/internal/clog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Ck8s struct {
	Clientset *kubernetes.Clientset
}

var ck8s Ck8s

func NewClient() error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		clog.Error(context.Background(), "error creating in-cluster config", err)
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		clog.Error(context.Background(), "error creating kubernetes clientset", err)
		return err
	}

	// set the client
	ck8s = Ck8s{
		Clientset: clientset,
	}

	return nil
}
