package providers

import (
	"fmt"
)

type KiamConigEnhancer struct {
	RoleName string		`yaml:"roleName"`
	ExternalId string	`yaml:"externalId"`
	AccountId string	`yaml:"accountId"`
	Partition string	`yaml:"partition"`
	KeySuffix string	`yaml:"keySuffix"`
}

func (enhancer *KiamConigEnhancer) Enhance(config *JobConfig) {
	var role = fmt.Sprintf("arn:%s:iam::%s:role/%s", enhancer.Partition, enhancer.AccountId, enhancer.RoleName)

	config.Annotations["iam.amazonaws.com/role"] = role
	config.Annotations["iam.amazonaws.com/external-id"] = enhancer.ExternalId
}

func (enhancer *KiamConigEnhancer) Key() string {
	if enhancer.KeySuffix != "" {
		return "aws-kiam-" + enhancer.KeySuffix
	}
	return "aws-kiam"
}
