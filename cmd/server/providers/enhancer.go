package providers

type JobConfig struct {
	Annotations map[string]string
	Labels map[string]string
}

type ConfigEnhancer interface {
	Enhance(config *JobConfig)
}
