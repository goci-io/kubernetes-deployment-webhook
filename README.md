# kubernetes-deployment-webhook

![Validate](https://github.com/goci-io/kubernetes-deployment-webhook/workflows/Validate/badge.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/goci-io/kubernetes-deployment-webhook/badge.svg?branch=master)](https://coveralls.io/github/goci-io/kubernetes-deployment-webhook?branch=master)

**Maintained by [@goci-io/prp-kubernetes](https://github.com/orgs/goci-io/teams/prp-kubernetes)**

### Purpose 

HTTP Server listening for Webhooks from Version Control Systems to deploy Kubernetes Jobs utilizing In-Cluster Auth to connect to the Kubernetes API and use In-Cluster Secrets. 
This Application can currently only validate Github Webhook Signatures and deploy from publicly available Repositories. Further restrictions initially apply by default to enhance security for example by disabling webhook from forked Repositories and limited to Releases or pushes to a references ending with `/master`. 

### Run

1. Build the Binary
- `make` (build within Docker, no Go required locally)
- `make bin/server` (Linux)
- `make bin/server/darwin` (MacOS).   
- Different GOOS: Run `GOOS=<GOOS> go build -o ./bin/webhook-server ./cmd/server` by your own

2. Configure Environment

See [Configure](https://github.com/goci-io/kubernetes-deployment-webhook/blob/master/README.md#configure) section on how to add additional Configuration files and configure the Webhook.

3. Run  
```
# Use Make to set required Env-Vars
# Defaults: WEBHOOK_SECRET=test, ORGANIZATION_WHITELIST=goci-io, disabled https
make run/local
```

You can also use our Docker Release:
```
docker run \
    -e CONFIG_DIR=/run/config
    -e WEBHOOK_SECRET=my-secret \
    -e ORGANIZATION_WHITELIST=org1,org2 \
    -v config:/run/config
    -it gocidocker/k8s-deployment-webhook:v0.1.0
```

You can find an example Request Payload [here](https://github.com/goci-io/kubernetes-deployment-webhook/blob/master/README.md#development).
By default the Application will serve its Endpoints via HTTPS (requires TLS configuration).

### Configure
The following Two Configuration Files are required:

##### [`repos.yaml`](https://github.com/goci-io/kubernetes-deployment-webhook/blob/master/config/repos.yaml)
Configure your Repositories and which Jobs to execute.

##### [`enhancers.yaml`](https://github.com/goci-io/kubernetes-deployment-webhook/blob/master/config/enhancers.yaml)
Enhancers are used to populate additional Fields into the Kubernetes Job Manifest.
You can read more about Enhancers within the [k8s package](https://github.com/goci-io/kubernetes-deployment-webhook/tree/master/cmd/kubernetes).

##### Webhook
- Webhook Secrets are required. The Server wont start without a Secret.
- Currently only one Webhook Secret is required per Installation
- Webhook Events need to be send to `POST /event`

##### `git pull` 
We can add an Init-Container to your Job automatically which downloads your Sources into a Temporary Directory.
This requires a Secret following the naming Convention of `<org>-<repo>-ssh`, providing `id_rsa` property containing a private [Deploy SSH Key](https://docs.github.com/en/developers/overview/managing-deploy-keys). Sources are mounted into `/run/workspace/checkout` with ReadOnly Access.

### Deploy

Using [goci-service-chart](https://github.com/goci-io/goci-service-chart) with the following **example config**: 

1. Create a Secret containing the `WEBHOOK_SECRET` environment variable  
2. Ensure following config environment variables are set as well: `ORGANIZATION_WHITELIST`  
3. Create ConfigMaps for repos.yaml and enhancers.yaml configuration  
4. Configure the following `values.yaml`:  
```yaml
port: 8443
fullnameOverride: k8s-deploy-webhook

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
  name: gocidocker/k8s-deploy-webhook
  tag: <LATEST_VERSION>

ingress:
  create: true
  tls: true
  hosts:
  - host: <YOUR_DOMAIN>
    paths:
    - /

extraVolumeMounts:
- name: tls
  mountPath: "/run/secrets/tls"

extraVolumes:
- name: tls
  secret:
    defaultMode: 0400
    secretName: k8s-deploy-webhook-tls
```

### Tests

Run Tests using `make tests`. This will run Tests in `server` and `kubernetes` packages.

### Development

To run this Application locally you need to build the Binary (see above) and run the App locally. To start the Application you can use `make run/local` which creates all necessary environment variables and default configurations (can be found [here](https://github.com/goci-io/kubernetes-deployment-webhook/tree/master/config/)). You will also need to specify the following Headers in your Requests:
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

#### Disable TLS
In case an external Provider is terminating TLS Connections you can force the Application to start a non-TLS Server by adding `FORCE_NON_TLS_SERVER=1` to your environment specification.

_This repository was created via [github-repository](https://github.com/goci-io/github-repository)._
