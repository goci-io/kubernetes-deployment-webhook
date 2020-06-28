package main

import (
	"os"
	"log"
	"errors"
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
		webhook, code, err := validateAndParseRequest(r, handler)

		if err != nil && code > 399 {
			failRequest(w, code, err)
		} else {
			succeedRequest(w, webhook, code)
		}
	}))

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}

func validateAndParseRequest(r *http.Request, handler *WebhookHandler) (*WebhookContext, int, error) {
	log.Print("Handling webhook request ...")

	webhook := &WebhookContext{}
	body, code, err := handler.validateRequest(r);
	if err != nil {
		return webhook, code, err
	}

	webhook, err = handler.parse(body, webhook)
	if err != nil {
		return webhook, code, err
	}

	if !handler.isEligible(webhook) {
		return webhook, http.StatusOK, errors.New("webhook is not eligible for processing")
	}

	return webhook, http.StatusAccepted, nil
}

func succeedRequest(w http.ResponseWriter, webhook *WebhookContext, code int) {
	log.Printf("Webhook request handled successfully: %v", webhook)
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
