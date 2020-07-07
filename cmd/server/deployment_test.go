package main

import (
	"errors"
	"strings"
	"testing"
	"github.com/goci-io/deployment-webhook/cmd/server/clients"
)

type KubernetesTestClient struct {
}

func (c *KubernetesTestClient) CreateJob(job *clients.DeploymentJob) error {
	if !strings.HasPrefix(job.Name, "goci-io-example-") {
		return errors.New("got invalid job name: " + job.Name)
	}

	return nil
}

func TestReleaseCallsKubernetesClientWithCorrectJobName(t *testing.T) {
	d := &Deployment{
		kubernetes: &KubernetesTestClient{},
	}

	ctx := &WebhookContext{
		Organization: "goci-io",
		Repository: &Repository{
			Name: "example",
		},
	}

	err := d.release(ctx)

	if err != nil {
		t.Error("deployment failed: " + err.Error())
	}
}

// @TODO add tests for enhancers
