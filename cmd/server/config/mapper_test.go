package config

import (
	"testing"
)

func TestMapperCreatesRepositoryConfigs(t *testing.T) {
	mapper := &DeploymentsConfig{}
	err := mapper.LoadAndParse("../../../config/example.yaml")

	if err != nil {
		t.Error("expected no error, got " + err.Error())
	}

	if cc := len(mapper.configs); cc != 1 {
		t.Errorf("expected one repository config got %d", cc)
	}

	repo := mapper.GetForRepo("goci-io", "example-repository")
	expected := &RepositoryConfig{
		Namespace: "default",
		Organization: "goci-io",
		Image: "repo/image:tag",
		Repository: "example-repository",
	}

	if !expected.equals(repo) {
		t.Errorf("expected %v, got %v", expected, repo)
	}
}
