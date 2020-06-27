package main

type Informer interface {
	inform(err error)
}

type KubernetesClient interface {
	createJob()
}

type Deployment struct {
	FailureInformer Informer
	SuccessInformer Informer
	Kubernetes KubernetesClient
}

func (deployment *Deployment) release(context *WebhookContext) {

}
