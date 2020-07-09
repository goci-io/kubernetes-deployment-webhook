package main

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var configs = make(map[string]RepositoryConfig)

type RepositoryConfig struct {
	Organization string		`yaml:"organization"`
	Repository string		`yaml:"repository"`
	ServiceAccount string	`yaml:"serviceAccount"`
	Namespace string		`yaml:"namespace"`
	Image string			`yaml:"image"`
	Providers []string		`yaml:"providers"`
}

func (config *RepositoryConfig) equals(other RepositoryConfig) bool {
	return config.Organization == other.Organization && 
		config.Repository == other.Repository &&
		config.Namespace == other.Namespace &&
		config.Image == other.Image
}

func LoadAndParseRepoConfig(path string) (map[string]RepositoryConfig, error) {
	configs := make(map[string]RepositoryConfig)
	yamlFile, err := ioutil.ReadFile(path)
    if err != nil {
        return configs, err
	}

	repos := []RepositoryConfig{}
	err = yaml.Unmarshal(yamlFile, &repos)
    if err != nil {
		return configs, err
	}

	for i := 0; i < len(repos); i++ {
		repoConfig := repos[i]
		key := fmt.Sprintf("%s/%s", repoConfig.Organization, repoConfig.Repository)

		configs[key] = repoConfig
	}

	return configs, nil
}
