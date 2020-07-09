package main

import (
	"testing"
)

func TestMapperCreatesRepositoryConfigs(t *testing.T) {
	repos, err := LoadAndParseRepoConfig("../../config/repos.yaml")

	if err != nil {
		t.Error("expected no error, got " + err.Error())
	}

	expected := &RepositoryConfig{
		Namespace: "default",
		Organization: "goci-io",
		Image: "repo/image:tag",
		Repository: "goci-repository-setup-example",
	}

	repo := repos["goci-io/goci-repository-setup-example"]
	if !expected.equals(repo) {
		t.Errorf("expected %v, got %v", expected, repo)
	}
}
