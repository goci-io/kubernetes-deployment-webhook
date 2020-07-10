package k8s

import (
	"errors"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	batchv1 "k8s.io/api/batch/v1"
)

type Enhancer interface {
	EnhanceJob(job *batchv1.Job)
	SetDefaults()
	Key() string
}

type ProviderConfig struct {
	Provider string				`yaml:"provider"`
	Config map[string]string	`yaml:",inline"`
}

func loadAndParseEnhancers(path string) ([]Enhancer, error) {
	configs := []ProviderConfig{}
	enhancers := []Enhancer{}

	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        return enhancers, err
	}

	err = yaml.Unmarshal(yamlFile, &configs)
    if err != nil {
		return enhancers, err
	}

    for i := range configs {
		provider := &configs[i]
		attributes, _ := yaml.Marshal(provider.Config)
		enhancer, err := unmarshalEnhancerAttributes(provider, attributes);
		if err != nil {
			return enhancers, err
		}

		enhancer.SetDefaults()
		enhancers = append(enhancers, enhancer)
	}

	return enhancers, nil
}

func unmarshalEnhancerAttributes(config *ProviderConfig, b []byte) (Enhancer, error) {
	switch config.Provider {
	case "aws-kiam":
		kiam := &KiamConigEnhancer{}
		err := yaml.Unmarshal(b, kiam)
		return kiam, err
	case "git-pull":
		gp := &PullGitSourcesEnhancer{}
		err := yaml.Unmarshal(b, gp)
		return gp, err
	default:
		return nil, errors.New("unknown provider " + config.Provider)
	}
}
