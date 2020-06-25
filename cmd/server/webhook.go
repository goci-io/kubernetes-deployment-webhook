package main

import (
	"fmt"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

const (
	jsonContentType = `application/json`
)

type Repository struct {
	Private bool `json:"private"`
}

type WebhookContext struct {
	Action string 		  `json:"action"`
	Repository Repository `json:"repository"`
	Organization string   `json:"organization"`
}

type WebhookHandler struct {
	Secret []byte
	OrganizationWhitelist []string
}

func (handler *WebhookHandler) isEligible(webhook *WebhookContext) bool {
	if webhook.Action != "published" {
		return false
	}

	if webhook.Repository.Private {
		return false
	}

	if findIndex(handler.OrganizationWhitelist, webhook.Organization) < 0 {
		return false
	}

	return true
}

func (handler *WebhookHandler) validateRequest(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("invalid method %s, only POST requests are allowed", r.Method)
	}

	if contentType := r.Header.Get("Content-Type"); contentType != jsonContentType {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("unsupported content type %s, only %s is supported", contentType, jsonContentType)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("could not read request body: %v", err)
	}

	delivery := r.Header.Get("x-github-delivery")
	signature := r.Header.Get("x-hub-signature")
	event := r.Header.Get("x-github-event")

	if len(signature) == 0 || len(event) == 0 || len(delivery) == 0 || !verifySignature(handler.Secret, signature, body) {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New("missing github event signature, webhook event, id or signature is invalid")
	}

	return body, nil
}

func (handler *WebhookHandler) parse(body []byte) (*WebhookContext, error) {
	webhook := &WebhookContext{}
	if len(body) == 0 {
		return webhook, errors.New("request body is empty")
	}

	err := json.Unmarshal([]byte(body), *webhook)
	if err != nil {
		return webhook, errors.New("invalid request, could not parse webhook object")
	}

	return webhook, nil
}

// https://gist.github.com/rjz/b51dc03061dbcff1c521
func verifySignature(secret []byte, signature string, body []byte) bool {
	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func findIndex(arr []string, search string) int {
    for i, n := range arr {
        if search == n {
            return i
        }
    }
    return -1
}
