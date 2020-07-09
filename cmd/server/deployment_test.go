package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/goci-io/deployment-webhook/cmd/kubernetes"
)

type KubernetesTestClient struct {
}

func (c *KubernetesTestClient) CreateJob(job *k8s.DeploymentJob) error {
	if !strings.HasPrefix(job.Name, "goci-io-example-") {
		return errors.New("got invalid job name: " + job.Name)
	}

	return nil
}

func TestReleaseCallsKubernetesClientWithCorrectJobName(t *testing.T) {
	d := &DeploymentsHandler{
		kubernetes: &KubernetesTestClient{},
	}

	ctx := &WebhookContext{
		Repository: &Repository{
			Organization: "goci-io",
			Name: "example",
		},
	}

	err := d.deploy(ctx)

	if err != nil {
		t.Error("deployment failed: " + err.Error())
	}
}

// @TODO add tests for enhancers
