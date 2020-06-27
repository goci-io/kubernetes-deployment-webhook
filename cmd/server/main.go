package main

import (
	"os"
	"log"
	"strings"
	"net/http"
	"path/filepath"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

func main() {
	handler := &WebhookHandler{
		Secret: []byte(os.Getenv("WEBHOOK_SECRET")),
		OrganizationWhitelist: strings.Split(os.Getenv("ORGANIZATION_WHITELIST"), ","),
	}

	if len(handler.Secret) == 0 {
		log.Fatal("missing required webhook sercret")
	}

	mux := http.NewServeMux()
	mux.Handle("/event", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateAndParseRequest(w, r, handler)
	}))

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}

func validateAndParseRequest(w http.ResponseWriter, r *http.Request, handler *WebhookHandler) {
	log.Print("Handling webhook request ...")

	body, code, err := handler.validateRequest(r);
	if err != nil {
		failRequest(w, code, err)
		return
	}

	webhook, err := handler.parse(body)
	if err != nil {
		failRequest(w, code, err)
		return
	}

	if !handler.isEligible(webhook) {
		succeedRequest(w, http.StatusOK)
		return
	}

	//deployment := &Deployment{Request: webhook}
	//defer deployment.release()

	succeedRequest(w, code)
}

func succeedRequest(w http.ResponseWriter, code int) {
	log.Print("Webhook request handled successfully")
	w.WriteHeader(code)
}

func failRequest(w http.ResponseWriter, code int, err error) {
	log.Printf("Error handling webhook request: %v", err)

	w.WriteHeader(code)
	_, writeErr := w.Write([]byte(err.Error()))

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}
