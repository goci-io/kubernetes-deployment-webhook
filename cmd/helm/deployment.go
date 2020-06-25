package main

type Deployment struct {
	Request *WebhookContext
}

func (deployment *Deployment) release() (bool, error) {
	return true, nil
}
