package config

import (
	"testing"
)

func TestMapperCreatesRepositoryConfigs(t *testing.T) {
	err := LoadAndParse("../../../config/repos.yaml")

	if err != nil {
		t.Error("expected no error, got " + err.Error())
	}

	repo := GetForRepo("goci-io", "goci-repository-setup-example")
	expected := &RepositoryConfig{
		Namespace: "default",
		Organization: "goci-io",
		Image: "repo/image:tag",
		Repository: "goci-repository-setup-example",
	}

	if !expected.equals(repo) {
		t.Errorf("expected %v, got %v", expected, repo)
	}
}
