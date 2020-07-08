package main

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	jsonContentType = `application/json`
)

type Release struct {
	TagName string `json:"tag_name"`
}

type Repository struct {
	Name string				`json:"name"`
	Fork bool				`json:"fork"`
	Private bool 			`json:"private"`
	Organization string    	`json:"organization"`
}

type WebhookContext struct {
	Action string 		   `json:"action"`
	Repository *Repository `json:"repository"`
	Reference string	   `json:"ref,omitempty"`
	Release *Release       `json:"release,omitempty"`
}

type WebhookHandler struct {
	secret []byte
	gitHost string
	kubernetes KubernetesClient
	organizationWhitelist []string
}

func (handler *WebhookHandler) handle(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling webhook request ...")

	webhook := &WebhookContext{}
	body, code, err := handler.validateRequest(r);
	if err != nil {
		failRequest(w, code, err)
		return
	}

	webhook, err = handler.parse(body, webhook)
	if err != nil {
		failRequest(w, http.StatusBadRequest, err)
		return
	}

	if !handler.isEligible(webhook) {
		succeedRequest(w, http.StatusOK)
		return 
	}

	deployment := &Deployment{
		kubernetes: handler.kubernetes,
	}

	err = deployment.release(webhook)
	if err != nil {
		failRequest(w, 500, err)
	} else {
		succeedRequest(w, code)
	}
}

func (handler *WebhookHandler) validateRequest(r *http.Request) ([]byte, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("invalid method %s, only POST requests are allowed", r.Method)
	}

	if contentType := r.Header.Get("Content-Type"); contentType != jsonContentType {
		return nil, http.StatusBadRequest, fmt.Errorf("unsupported content type %s, only %s is supported", contentType, jsonContentType)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("could not read request body: %v", err)
	}

	signature := r.Header.Get("x-hub-signature")
	event := r.Header.Get("x-github-event")

	if len(signature) == 0 || len(event) == 0 {
		return body, http.StatusBadRequest, errors.New("missing github event signature, webhook event, id or signature is invalid")
	}

	if !verifySignature(handler.secret, signature, body) {
		return body, http.StatusBadRequest, errors.New("invalid webhook signature")
	}

	return body, http.StatusAccepted, nil
}

func (handler *WebhookHandler) parse(body []byte, into *WebhookContext) (*WebhookContext, error) {
	if len(body) == 0 {
		return into, errors.New("request body is empty")
	}

	err := json.Unmarshal(body, into)
	if err != nil {
		return into, errors.New("invalid request, could not parse webhook object: " + err.Error())
	}

	return into, nil
}

func (handler *WebhookHandler) isEligible(webhook *WebhookContext) bool {
	isRelease := webhook.Action == "published" && webhook.Release != nil
	isMasterMerge := len(webhook.Reference) > 0 && strings.HasSuffix(webhook.Reference, "/master")

	if !isRelease && !isMasterMerge {
		return false
	}

	if webhook.Repository == nil || webhook.Repository.Private || webhook.Repository.Fork {
		return false
	}

	if !contains(handler.organizationWhitelist, webhook.Repository.Organization) {
		return false
	}

	return true
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
