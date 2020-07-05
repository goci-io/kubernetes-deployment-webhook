package main

import (
	"github.com/goci-io/deployment-webhook/cmd/server/providers"
	"github.com/goci-io/deployment-webhook/cmd/server/clients"
	"github.com/goci-io/deployment-webhook/cmd/server/config"
)

type Informer interface {
	inform(err error)
}

type Deployment struct {
	FailureInformer Informer
	SuccessInformer Informer
	Configs config.DeploymentsConfig
	Kubernetes clients.KubernetesClient
	Enhancer []providers.ConfigEnhancer
}

func (deployment *Deployment) release(context *WebhookContext) {
	config := deployment.Configs.GetForRepo(context.Organization, context.Repository.Name)
	job := &clients.DeploymentJob{}
	copyConfigInto(&config, job)

	enhanced := &providers.JobConfig{}
	for i := 0; i < len(deployment.Enhancer); i++ {
		enhancer := deployment.Enhancer[i]
		enhancer.Enhance(enhanced)
	}
}

func copyConfigInto(config *config.RepositoryConfig, into *clients.DeploymentJob) {
	into.Image = config.Image
	into.Namespace = config.Namespace
	into.ServiceAccount = config.ServiceAccount
	into.Labels = make(map[string]string)
	into.Labels["repository"] = config.Repository
	into.Labels["organization"] = config.Organization
}
