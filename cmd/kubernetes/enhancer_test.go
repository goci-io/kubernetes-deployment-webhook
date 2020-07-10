package k8s

import (
	"testing"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnhancerLoadAndParseCreatesEnhancers(t *testing.T) {
	enhancers, err := loadAndParseEnhancers("../../config/enhancers.yaml")

	if err != nil {
		t.Error("error loading providers config: " + err.Error())
	}

	if len(enhancers) != 1 {
		t.Error("expected exactly one example kiam provider to be configured")
	}

	kiam := enhancers[0].(*KiamConigEnhancer)
	if kiam.KeySuffix != "goci-app" || kiam.Key() != "aws-kiam-goci-app" {
		t.Error("key suffix not correctly mapped. got: " + kiam.KeySuffix)
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: make(map[string]string),
		},
	}

	kiam.EnhanceJob(job)
	expectedRole := "arn:aws:iam::123456789012:role/goci-build-app-role"

	if job.Annotations["iam.amazonaws.com/role"] != expectedRole {
		t.Errorf("expected role %s, got %s", expectedRole, job.Annotations["iam.amazonaws.com/role"])
	}
}
