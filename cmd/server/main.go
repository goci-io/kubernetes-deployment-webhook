package main

import (
	"os"
	"log"
	"strings"
	"net/http"
	"path/filepath"
	"github.com/goci-io/deployment-webhook/cmd/server/config"
	"github.com/goci-io/deployment-webhook/cmd/server/clients"
	"github.com/goci-io/deployment-webhook/cmd/server/providers"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

func main() {
	k8sClient := &clients.KubernetesClient{}
	k8sClient.Init()

	webhook := &WebhookHandler{
		kubernetes: k8sClient,
		gitHost: getEnv("GIT_HOST", "github.com"),
		secret: []byte(os.Getenv("WEBHOOK_SECRET")),
		organizationWhitelist: strings.Split(os.Getenv("ORGANIZATION_WHITELIST"), ","),
	}

	if len(webhook.secret) == 0 {
		log.Fatal("missing required webhook sercret")
	}

	initRandom()

	config.LoadAndParse(getEnv("REPO_CONFIG_FILE", "/run/config/repos.yaml"))
	providers.LoadAndParse(getEnv("PROVIDERS_CONFIG_FILE", "/run/config/providers.yaml"))

	mux := http.NewServeMux()
	mux.Handle("/event", http.HandlerFunc(webhook.handle))

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
