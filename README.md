# kubernetes-deployment-webhook

![Validate](https://github.com/goci-io/kubernetes-deployment-webhook/workflows/Validate/badge.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/goci-io/kubernetes-deployment-webhook/badge.svg?branch=master)](https://coveralls.io/github/goci-io/kubernetes-deployment-webhook?branch=master)

**Maintained by [@goci-io/prp-kubernetes](https://github.com/orgs/goci-io/teams/prp-kubernetes)**

### Purpose 

HTTP Server listening for Webhooks from Version Control Systems to deploy Kubernetes Jobs utilizing In-Cluster Auth to connect to the Kubernetes API. 
This Application can currently only validate Github Webhook Signatures and deploy from publicly available Repositories. Further restrictions initially apply by default to enhance security for example by disabling webhook from forked Repositories and limited to Releases or pushes to a references ending with `/master`. 

At goci.io we use this Webhook Server for our external Provider-Integrations to support their own Release-Cycles and further configuration for Deployments.

### Run

1. Build the Binary
`make image/server` or `make image/server/darwin`.   
1.1 In case you are running on a different GOOS than Linux or Darwin you need to use `GOOS=<GOOS> go build -o ./bin/webhook-server ./cmd/server` by your own.

2. Configure Environment
```
export GIT_HOST=github.com (default)
export WEBHOOK_SECRET=my-secret (required)
export ORGANIZATION_WHITELIST=org1,org2 (default: none)
```
3. Run
```
./bin/webhook-server
```

You can also use our Docker Release:
```
docker run \
    -e WEBHOOK_SECRET=my-secret \
    -e ORGANIZATION_WHITELIST=org1,org2 \
    -it gocidocker/k8s-deployment-webhook:v0.1.0
```

### Deploy

Using [goci-service-chart](https://github.com/goci-io/goci-service-chart) with the following **example config**: 

1. Create a Secret containing the `WEBHOOK_SECRET` environment variable
2. Ensure following config environment variables are set as well: `ORGANIZATION_WHITELIST`
2. Configure the following `values.yaml`:  
```yaml
port: 8443

configMap:
  data:
    ORGANIZATION_WHITELIST: org1,org2

envFrom:
- secretRef:
    name: <YOUR_SECRET_NAME>

# Allow deploying Jobs with different SAs
rbac:
  appUser:
    create: true # TBD with goci-service-chart
  additionalRules:
  - apiGroups: ["batch/v1"]
    resources: ["jobs"]
    verbs: ["create"]

image:
  name: gocidocker/k8s-deployment-webhook
  tag: v0.1.0

# Mount TLS Secrets
extraVolumeMounts:
- name: tls
  mountPath: "/run/secrets/tls"

extraVolumes:
- name: tls
  secret:
    defaultMode: 0600
    secretName: k8s-deployment-webhook-tls
```

### Tests

Run Tests using `make tests`.

To run this Application locally you need to build the Binary (see above) and run the App locally. To start the Application you can use `make run/local` which creates all necessary environment variables and default configurations (can be found [here](config/)). You will also need to specify the following Headers in your Requests:
```
x-hub-signature: sha1=37df8bce63ad1e1acc699d9575afabc3de2ff9ac
x-github-event: push
```

<details><summary>Example Webhook Payload</summary>

```json
{
  "ref": "refs/heads/master",
  "repository": {
    "name": "goci-repository-setup-example",
    "full_name": "goci-io/goci-repository-setup-example",
    "private": false,
    "owner": {
      "name": "goci-io",
      "email": "support@goci.io",
      "login": "goci-io",
      "url": "https://api.github.com/users/goci-io",
      "type": "Organization"
    },
    "fork": false,
    "url": "https://github.com/goci-io/goci-repository-setup-example",
    "git_url": "git://github.com/goci-io/goci-repository-setup-example.git",
    "ssh_url": "git@github.com:goci-io/goci-repository-setup-example.git",
    "clone_url": "https://github.com/goci-io/goci-repository-setup-example.git",
    "default_branch": "master",
    "master_branch": "master",
    "organization": "goci-io"
  },
  "pusher": {
    "name": "etwillbefine",
    "email": "etwillbefine@users.noreply.github.com"
  },
  "organization": {
    "login": "goci-io",
    "url": "https://api.github.com/orgs/goci-io",
    "repos_url": "https://api.github.com/orgs/goci-io/repos",
    "events_url": "https://api.github.com/orgs/goci-io/events",
    "hooks_url": "https://api.github.com/orgs/goci-io/hooks"
  }
}
```
</details>

_This repository was created via [github-repository](https://github.com/goci-io/github-repository)._
