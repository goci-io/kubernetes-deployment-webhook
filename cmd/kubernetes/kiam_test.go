package k8s

import (
	"testing"
)

func TestKiamEnhancerAppendsIAMAnnotations(t *testing.T) {
	job := &DeploymentJob{Annotations: make(map[string]string)}
	enhancer := &KiamConigEnhancer{
		Partition: "aws",
		RoleName: "example",
		AccountId: "12345678912",
		ExternalId: "external-id",
	}

	enhancer.Enhance(job)
	expectedRole := "arn:aws:iam::12345678912:role/example"

	if job.Annotations["iam.amazonaws.com/role"] != expectedRole {
		t.Errorf("expected role %s, got %s", expectedRole, job.Annotations["iam.amazonaws.com/role"])
	}

	if job.Annotations["iam.amazonaws.com/external-id"] != enhancer.ExternalId {
		t.Errorf("expected role %s, got %s", enhancer.ExternalId, job.Annotations["iam.amazonaws.com/role"])
	}
}
