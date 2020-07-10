# `k8s` package

### Contains
- Kubernetes Client
- Enhancer for Kubernetes Manifests
- Loads YAML Configuration for Enhancers
- Able to detect In-Cluster Auth automatically and fallback to `kubecfg`

### Usage

```go
import (
	"github.com/goci-io/deployment-webhook/cmd/kubernetes"
)

func main() {
    k8sClient := &k8s.Client{}
    k8sClient.Init("/path/to/config.yaml")
    
    d := &k8s.DeploymentJob{
        // ...

        // Activate/Use Enhancers
        Enhancers: []string{"git-pull", "aws-kiam"},
        Data: {
            // pass custom data
        },
    }

    k8sClient.CreateJob(d)
}
```

### Configuration
```yaml
- externalId: AvoidConfusedDeputyProblem
  roleName: goci-build-app-role
  accountId: 123456789012
  provider: aws-kiam
- ...
```

### Enhancers

The following Providers are available:

1. [`aws-kiam`](kiam.go)
1. [`git-pull`](git_pull.go)
