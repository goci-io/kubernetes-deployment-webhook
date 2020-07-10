package k8s

import (
	"testing"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKiamEnhancerAppendsIAMAnnotations(t *testing.T) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
		},
	}

	enhancer := &KiamConigEnhancer{
		Partition: "aws",
		RoleName: "example",
		AccountId: "12345678912",
		ExternalId: "external-id",
	}

	enhancer.EnhanceJob(job)
	expectedRole := "arn:aws:iam::12345678912:role/example"

	if job.Annotations["iam.amazonaws.com/role"] != expectedRole {
		t.Errorf("expected role %s, got %s", expectedRole, job.Annotations["iam.amazonaws.com/role"])
	}

	if job.Annotations["iam.amazonaws.com/external-id"] != enhancer.ExternalId {
		t.Errorf("expected role %s, got %s", enhancer.ExternalId, job.Annotations["iam.amazonaws.com/role"])
	}
}
