package k8s

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PullGitSourcesEnhancer struct {
	GitHost string `yaml:"host"`
	GitUser string `yaml:"user"`
	KeySuffix string `yaml:"keySuffix,omitempty"`
}

type PullGitSourcesData interface {
	JobData
	Repository() string
	Organization() string
}

func (enhancer *PullGitSourcesEnhancer) EnhanceJob(job *batchv1.Job, d JobData) {
	data := d.(PullGitSourcesData)

	if job.Spec.Template.Spec.InitContainers == nil {
		job.Spec.Template.Spec.InitContainers = []corev1.Container{}
	}
	if job.Spec.Template.Spec.Volumes == nil {
		job.Spec.Template.Spec.Volumes = []corev1.Volume{}
	}
	if job.Spec.Template.Spec.Containers[0].VolumeMounts == nil {
		job.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{}
	}

	job.Spec.Template.Spec.Containers[0].VolumeMounts = append(job.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name: "sources",
		ReadOnly: false,
		MountPath: "/run/workspace",
	})

	pullCmd := fmt.Sprintf("git clone %s@%s:%s/%s.git",
		enhancer.GitUser, enhancer.GitHost, data.Organization(), data.Repository())

	job.Spec.Template.Spec.InitContainers = append(job.Spec.Template.Spec.InitContainers, corev1.Container{
		Name: "pull-sources",
		Image: "gocidocker/k8s-deploy-alpine:0.1.1",
		Command: []string{"/bin/bash"},
		Args: []string{"-c", pullCmd},
		Env: []corev1.EnvVar{
			{
				Name: "GIT_SSH_COMMAND",
				Value: "ssh -i /run/secrets/git/id_rsa",
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name: "sources",
				ReadOnly: true,
				MountPath: "/run/workspace/checkout",
			},
			{
				Name: "git-ssh",
				ReadOnly: true,
				MountPath: "/run/secrets/git",
			},
		},
	})

	job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes, 
		corev1.Volume{
			Name: "sources",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					SizeLimit: resource.NewQuantity(4 * 1000*1000*1000, resource.DecimalSI),
				},
			},
		},
		corev1.Volume{
			Name: "git-ssh",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					Optional: createBool(false),
					DefaultMode: create32(256),
					SecretName: fmt.Sprintf("%s-%s-ssh", data.Organization(), data.Repository()),
				},
			},
		},
	)
}

func (enhancer *PullGitSourcesEnhancer) Key() string {
	if enhancer.KeySuffix != "" {
		return "git-pull-" + enhancer.KeySuffix
	}
	return "git-pull"
}

func (enhancer *PullGitSourcesEnhancer) SetDefaults() {
	if enhancer.GitHost == "" {
		enhancer.GitHost = "github.com"
	}
	if enhancer.GitUser == "" {
		enhancer.GitUser = "git"
	}
}
