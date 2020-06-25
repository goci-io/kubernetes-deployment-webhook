package main

import (
	"strings"
	"testing"
	"net/http"
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
	if err == nil || err.Error() != "missing github event signature, webhook event, id or signature is invalid" {
		t.Error("expected invalid request content error")
	}

	req.Header.Add("x-github-signature", "sig")
	if err == nil || err.Error() != "missing github event signature, webhook event, id or signature is invalid" {
		t.Error("expected invalid request content error")
	}
}
