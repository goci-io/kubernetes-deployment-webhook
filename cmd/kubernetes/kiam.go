package k8s

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
)

type KiamConigEnhancer struct {
	RoleName string		`yaml:"roleName"`
	ExternalId string	`yaml:"externalId"`
	AccountId string	`yaml:"accountId"`
	Partition string	`yaml:"partition"`
	KeySuffix string	`yaml:"keySuffix,omitempty"`
}

func (enhancer *KiamConigEnhancer) EnhanceJob(job *batchv1.Job, data JobData) {
	var role = fmt.Sprintf("arn:%s:iam::%s:role/%s", enhancer.Partition, enhancer.AccountId, enhancer.RoleName)

	job.ObjectMeta.Annotations["iam.amazonaws.com/role"] = role
	job.ObjectMeta.Annotations["iam.amazonaws.com/external-id"] = enhancer.ExternalId
}

func (enhancer *KiamConigEnhancer) Key() string {
	if enhancer.KeySuffix != "" {
		return "aws-kiam-" + enhancer.KeySuffix
	}
	return "aws-kiam"
}

func (enhancer *KiamConigEnhancer) SetDefaults() {
	if enhancer.Partition == "" {
		enhancer.Partition = "aws"
	}
}
