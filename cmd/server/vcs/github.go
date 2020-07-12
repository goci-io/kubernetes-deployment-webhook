package vcs // import "github.com/goci-io/deployment-webhook/cmd/vcs"

import (
	"net/http"
)

type GithubProvider struct {
}

func (provider *GithubProvider) Signature(r *http.Request) string {
	return r.Header.Get("x-hub-signature")
}

func (provider *GithubProvider) Event(r *http.Request) string {
	return r.Header.Get("x-github-event")
}
