package config

import (
	"testing"
)

func TestMapperCreatesRepositoryConfigs(t *testing.T) {
	err := LoadAndParse("../../../config/example.yaml")

	if err != nil {
		t.Error("expected no error, got " + err.Error())
	}

	repo := GetForRepo("goci-io", "example-repository")
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
