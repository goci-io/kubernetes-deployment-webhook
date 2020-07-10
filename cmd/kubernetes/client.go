package k8s

import (
	"os"
	"flag"
	"strings"
	"path/filepath"

	//"k8s.io/apimachinery/pkg/api/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	batchv1types "k8s.io/client-go/kubernetes/typed/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// used within Enhancers for example
type JobData interface {}

type DeploymentJob struct {
	Name string
	TTL  int32
	Data JobData
	Image string
	Namespace string
	Enhancers []string
	ServiceAccount string
	Annotations map[string]string
	Labels map[string]string
	SecretEnvName string
}

type Client struct {
	enhancers []Enhancer
	batchV1 batchv1types.BatchV1Interface
}

func (client *Client) Init(enhancerConfigPath string) error {
	var config *rest.Config
	var err error

	if inClusterAuthPossible() {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	if err != nil {
		return err
	}

	clientsets, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	client.batchV1 = clientsets.BatchV1()
	client.enhancers, err = loadAndParseEnhancers(enhancerConfigPath)
	return err
}

func (client *Client) CreateJob(job *DeploymentJob) error {
	name := strings.ToLower(job.Name)

	manifest := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: job.Labels,
			Annotations: job.Annotations,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: create32(1),
			ActiveDeadlineSeconds: create64(10800),
			TTLSecondsAfterFinished: &job.TTL,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: job.Labels,
					Annotations: job.Annotations,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					ServiceAccountName: job.ServiceAccount,
					TerminationGracePeriodSeconds: create64(100),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: createBool(true),
						RunAsUser: create64(1000),
						FSGroup: create64(1000),
					},
					Containers: []corev1.Container{
						{
							Name: "job",
							Image: job.Image,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("300m"),
									corev1.ResourceMemory: resource.MustParse("156Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("300m"),
									corev1.ResourceMemory: resource.MustParse("156Mi"),
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: job.SecretEnvName,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i := 0; i < len(client.enhancers); i++ {
		enhancer := client.enhancers[i]

		if contains(job.Enhancers, enhancer.Key()) {
			enhancer.EnhanceJob(manifest, job.Data)
		}
	}

	_, err := client.batchV1.Jobs(job.Namespace).Create(manifest)
	return err
}

func createBool(x bool) *bool {
	return &x
}

func create32(x int32) *int32 {
    return &x
}

func create64(x int64) *int64 {
    return &x
}

func inClusterAuthPossible() bool {
	fi, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return os.Getenv("KUBERNETES_SERVICE_HOST") != "" &&
		os.Getenv("KUBERNETES_SERVICE_PORT") != "" &&
		err == nil && !fi.IsDir()
}

func contains(arr []string, search string) bool {
    for _, n := range arr {
        if search == n {
            return true
        }
    }
    return false
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
