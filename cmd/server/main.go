package main

import (
	"os"
	"log"
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

	body, err := handler.validateRequest(w, r);
	if err != nil {
		failRequest(w, err)
		return
	}

	webhook, err := handler.parse(body)
	if err != nil {
		failRequest(w, err)
		return
	}

	if !webhook.isEligible() {
		succeedRequest(w)
		return
	}

	//deployment := &Deployment{Request: webhook}
	//defer deployment.release()

	succeedRequest(w)
}

func succeedRequest(w http.ResponseWriter) {
	log.Print("Webhook request handled successfully")
	w.WriteHeader(http.StatusAccepted)
}

func failRequest(w http.ResponseWriter, err error) {
	log.Printf("Error handling webhook request: %v", err)

	_, writeErr := w.Write([]byte(err.Error()))

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}
