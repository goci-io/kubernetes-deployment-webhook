package main

import (
	"fmt"
	"log"
	"github.com/goci-io/deployment-webhook/cmd/server/providers"
	"github.com/goci-io/deployment-webhook/cmd/server/clients"
	"github.com/goci-io/deployment-webhook/cmd/server/config"
)

type MergeableMap map[string]string

type Informer interface {
	inform(err error)
}

type KubernetesClient interface {
	CreateJob(job *clients.DeploymentJob) error
}

type Deployment struct {
	failureInformer Informer
	successInformer Informer
	kubernetes KubernetesClient
	enhancers []providers.ConfigEnhancer
}

func (d *Deployment) release(context *WebhookContext) error {
	config := config.GetForRepo(context.Repository.Organization, context.Repository.Name)
	jobName := fmt.Sprintf("%s-%s-%s", context.Repository.Organization, context.Repository.Name, randStringBytes(6))
	job := &clients.DeploymentJob{Name: jobName}
	copyConfigInto(config, job)

	pd := &providers.JobConfig{
		Labels: make(MergeableMap),
		Annotations: make(MergeableMap),
	}

	for i := 0; i < len(d.enhancers); i++ {
		enhancer := d.enhancers[i]

		if contains(config.Providers, enhancer.Key()) {
			enhancer.Enhance(pd)
		}
	}

	mergeMap(job.Labels, pd.Labels);
	mergeMap(job.Annotations, pd.Annotations);

	log.Printf("deploying job %v", d.kubernetes.CreateJob(job))
	return d.kubernetes.CreateJob(job)
}

func mergeMap(target MergeableMap, merge MergeableMap) {
	for k, v := range merge {
		target[k] = v
	}
}

func copyConfigInto(config config.RepositoryConfig, into *clients.DeploymentJob) {
	into.Image = config.Image
	into.Namespace = config.Namespace
	into.ServiceAccount = config.ServiceAccount
	into.Labels = make(map[string]string)
	into.Labels["repository"] = config.Repository
	into.Labels["organization"] = config.Organization
}
