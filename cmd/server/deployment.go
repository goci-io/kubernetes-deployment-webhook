package main

import (
	"fmt"
	"github.com/goci-io/deployment-webhook/cmd/kubernetes"
)

type MergeableMap map[string]string

type Informer interface {
	inform(err error)
}

type KubernetesClient interface {
	CreateJob(job *k8s.DeploymentJob) error
}

type DeploymentsHandler struct {
	failureInformer Informer
	successInformer Informer
	enhancers []k8s.Enhancer
	kubernetes KubernetesClient
	configs map[string]RepositoryConfig
}

func (d *DeploymentsHandler) deploy(context *WebhookContext) error {
	jobName := fmt.Sprintf("%s-%s-%s", context.Repository.Organization, context.Repository.Name, randStringBytes(6))
	configName := fmt.Sprintf("%s-%s", context.Repository.Organization, context.Repository.Name)
	config := d.configs[configName]

	job := &k8s.DeploymentJob{
		Name: jobName,
		SecretEnvName: configName,
	}

	for i := 0; i < len(d.enhancers); i++ {
		enhancer := d.enhancers[i]

		if contains(config.Providers, enhancer.Key()) {
			enhancer.Enhance(job)
		}
	}

	copyConfigInto(config, job)
	return d.kubernetes.CreateJob(job)
}

func copyConfigInto(config RepositoryConfig, into *k8s.DeploymentJob) {
	into.Image = config.Image
	into.Namespace = config.Namespace
	into.ServiceAccount = config.ServiceAccount
	into.Labels = make(map[string]string)
	into.Labels["repository"] = config.Repository
	into.Labels["organization"] = config.Organization
}
