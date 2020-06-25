package main

import (
	"strings"
	"testing"
	"net/http"
	"encoding/hex"
	httptest "net/http/httptest"
)

func TestWebhookValidateRequestRejectsNonHttpPost(t *testing.T) {
	reader := strings.NewReader("request payload")
	req := httptest.NewRequest(http.MethodGet, "/", reader)
	writer := httptest.NewRecorder()

	handler := &WebhookHandler{}
	_, err := handler.validateRequest(writer, req)

	if err == nil || err.Error() != "invalid method GET, only POST requests are allowed" {
		t.Error("expected invalid http method error")
	}
}

func TestWebhookValidateRequestRejectsNonJsonContentType(t *testing.T) {
	reader := strings.NewReader("request payload")
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Content-Type", "application/yaml")

	handler := &WebhookHandler{}
	_, err := handler.validateRequest(httptest.NewRecorder(), req)

	if err == nil || err.Error() != "unsupported content type application/yaml, only application/json is supported" {
		t.Error("expected invalid content-type error")
	}
}

func TestWebhookValidateRequestRejectsMissingSignature(t *testing.T) {
	reader := strings.NewReader("request payload")
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Content-Type", "application/json")

	handler := &WebhookHandler{}
	_, err := handler.validateRequest(httptest.NewRecorder(), req)

	if err == nil || err.Error() != "missing github event signature, webhook event, id or signature is invalid" {
		t.Error("expected invalid request content error")
	}

	req.Header.Add("x-github-delivery", "x")
	req.Header.Add("x-github-event", "y")
	_, err = handler.validateRequest(httptest.NewRecorder(), req)
	if err == nil || err.Error() != "missing github event signature, webhook event, id or signature is invalid" {
		t.Error("expected invalid request content error")
	}

	req.Header.Add("x-hub-signature", "sig")
	_, err = handler.validateRequest(httptest.NewRecorder(), req)
	if err == nil || err.Error() != "missing github event signature, webhook event, id or signature is invalid" {
		t.Error("expected invalid request content error")
	}
}

func TestWebhookValidateRequestSucceeds(t *testing.T) {
	payload := "request payload"
	secret := "secret"
	signatureHmac := signBody([]byte(secret), []byte(payload))
	signature := make([]byte, hex.EncodedLen(len(signatureHmac)))
	hex.Encode(signature, signatureHmac)

	reader := strings.NewReader(payload)
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-github-event", "ping")
	req.Header.Add("x-hub-signature", "sha1=" + string(signature))
	req.Header.Add("x-github-delivery", "f155be00-b701-11ea-97da-4470faac803a")

	handler := &WebhookHandler{
		Secret: []byte("secret"),
	}

	body, err := handler.validateRequest(httptest.NewRecorder(), req)
	if err != nil {
		t.Error("expected valid request to succeed: " + err.Error())
	}

	parsedBody := string(body[:len(body)])
	if parsedBody != payload {
		t.Error("wrong content received")
	}
}

func TestWebhookIsEligibleForNonWhitelistedOrgFails(t *testing.T) {
	handler := &WebhookHandler{
		OrganizationWhitelist: []string{"goci-io", "goci-io-dev"},
	}

	webhook := &WebhookContext{
		Organization: "another-org",
	}

	eligible := handler.isEligible(webhook)

	if eligible {
		t.Error("expected another-org to be ineligible as its not a whitelisted org")
	}
}
