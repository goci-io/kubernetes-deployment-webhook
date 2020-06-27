package main

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClient struct {
	BatchV1 *batchv1.BatchV1Client
}

func (client *KubernetesClient) init() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientsets, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	client.BatchV1 = clientsets.BatchV1()
	return nil
}

func (client *KubernetesClient) createJob() error {
	return nil
}
