package main // import "github.com/goci-io/deployment-webhook/cmd/server"

import (
	"os"
	"log"
	"strings"
	"net/http"
	"path/filepath"
	"github.com/goci-io/deployment-webhook/cmd/kubernetes"
	"github.com/goci-io/deployment-webhook/cmd/server/vcs"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

func main() {
	initRandom()
	webhook := createHandler()

	mux := http.NewServeMux()
	mux.Handle("/event", http.HandlerFunc(webhook.handle))

	server := &http.Server{
		Addr:    ":9443",
		Handler: mux,
	}

	if getEnv("FORCE_NON_TLS_SERVER", "0") == "1" {
		log.Print("warn: using non-https server")
		log.Fatal(server.ListenAndServe())
	} else {
		certPath := filepath.Join(tlsDir, tlsCertFile)
		keyPath := filepath.Join(tlsDir, tlsKeyFile)
		log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
	}
}

func createHandler() *WebhookHandler {
	configPath := getEnv("CONFIG_DIR", "/run/config")
	enhancersPath := configPath + "/enhancers.yaml"
	reposPath := configPath + "/repos.yaml"

	configs, err := LoadAndParseRepoConfig(reposPath)
	if err != nil {
		log.Fatal("failed loading repository configuration: " + err.Error())
	}

	k8sClient := &k8s.Client{}
	if err = k8sClient.Init(enhancersPath); err != nil {
		log.Fatal("error creating k8s client: " + err.Error())
	}

	deployments := &DeploymentsHandler{
		configs: configs,
		kubernetes: k8sClient,
	}

	webhook := &WebhookHandler{
		deployments: deployments,
		vcsClient: &vcs.GithubProvider{},
		gitHost: getEnv("GIT_HOST", "github.com"),
		secret: []byte(os.Getenv("WEBHOOK_SECRET")),
		organizationWhitelist: strings.Split(os.Getenv("ORGANIZATION_WHITELIST"), ","),
	}

	if len(webhook.secret) < 1 {
		log.Fatal("missing required webhook sercret")
	}

	return webhook
}
