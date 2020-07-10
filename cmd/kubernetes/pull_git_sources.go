package k8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	batchv1 "k8s.io/api/batch/v1"
)

type PullGitSourcesEnhancer struct {
	GitHost string `yaml:"host"`
	GitUser string `yaml:"user"`
}

func (enhancer *PullGitSourcesEnhancer) EnhanceJob(job *batchv1.Job) {
	if job.Spec.Template.Spec.InitContainers == nil {
		job.Spec.Template.Spec.InitContainers = []corev1.Container{}
	}

	pullCmd := fmt.Sprintf("git clone %s@%s:%s/%s.git",
		enhancer.GitUser, enhancer.GitHost, job.Data.Organization, job.Data.Repository)

	append(job.Spec.Template.Spec.InitContainers, {
		Name: "pull-sources",
		Image: "gocidocker/k8s-deploy-alpine:0.1.0",
		Command: []string{"/bin/bash", "-c", pullCmd},
		Env: []corev1.EnvVar{
			{
				Name: "GIT_SSH_COMMAND",
				Value: "ssh -i /run/secrets/git/id_rsa",
			},
		},
		Volume: 
	})
}

func (enhancer *PullGitSourcesEnhancer) Key() string {
	if enhancer.KeySuffix != "" {
		return "git-pull-" + enhancer.KeySuffix
	}
	return "git-pull"
}

func (enhancer *PullGitSourcesEnhancer) SetDefaults() string {
	if enhancer.GitHost == "" {
		enhancer.GitHost = "github.com"
	}
	if enhancer.GitUser == "" {
		enhancer.GitUser = "git"
	}
}
