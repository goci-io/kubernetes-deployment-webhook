package k8s

import (
	"errors"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Enhancer interface {
	Enhance(config *DeploymentJob)
	Key() string
}

type ProviderConfig struct {
	Provider string				`yaml:"provider"`
	Config map[string]string	`yaml:",inline"`
}

func LoadAndParseEnhancers(path string) ([]Enhancer, error) {
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
	default:
		return nil, errors.New("unknown provider " + config.Provider)
	}
}
