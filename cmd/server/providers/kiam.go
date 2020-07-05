package providers

import (
	"fmt"
)

type KiamConigEnhancer struct {
	RoleName 	string
	ExternalId 	string
	AccountId	string
	Partition	string
}

func (enhancer *KiamConigEnhancer) Enhance(config *JobConfig) {
	var role = fmt.Sprintf("arn:%s:iam::%s:role/%s", enhancer.Partition, enhancer.AccountId, enhancer.RoleName)

	config.Annotations["iam.amazonaws.com/role"] = role
	config.Annotations["iam.amazonaws.com/external-id"] = enhancer.ExternalId
}
