package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type RepositoryConfig struct {
	Organization	string `yaml:"organization"`
	Repository		string `yaml:"repository"`
	ServiceAccount	string `yaml:"serviceAccount"`
	Namespace		string `yaml:"namespace"`
	Image			string `yaml:"image"`
}

func (config *RepositoryConfig) Equals(other RepositoryConfig) bool {
	return config.Organization == other.Organization && 
		config.Repository == other.Repository &&
		config.Namespace == other.Namespace &&
		config.Image == other.Image
}

type DeploymentsConfig struct {
	Configs map[string]RepositoryConfig
}

func (config *DeploymentsConfig) LoadAndParse(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        return err
	}

	configs := []RepositoryConfig{}
	err = yaml.Unmarshal(yamlFile, &configs)
    if err != nil {
		return err
	}
	
	config.Configs = make(map[string]RepositoryConfig)
	for i := 0; i < len(configs); i++ {
		repoConfig := configs[i]
		key := fmt.Sprintf("%s/%s", repoConfig.Organization, repoConfig.Repository)
		config.Configs[key] = repoConfig
	}

	return nil
}

func (config *DeploymentsConfig) GetForRepo(organization string, repository string) RepositoryConfig {
	key := fmt.Sprintf("%s/%s", organization, repository)
	return config.Configs[key]
}
