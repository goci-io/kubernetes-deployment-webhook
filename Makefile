.DEFAULT_GOAL := image/server

export CONFIG_DIR ?= ./config
export WEBHOOK_SECRET ?= test
export FORCE_NON_TLS_SERVER ?= 1
export ORGANIZATION_WHITELIST ?= goci-io

image/server:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

image/server/darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

image/docker: image/server
	docker build -t kubernetes-deployment-webhook .

run:
	go build

tests:
	go test github.com/goci-io/deployment-webhook/cmd/server/...
	go test github.com/goci-io/deployment-webhook/cmd/kubernetes/...

coverage:
	go test -v -coverprofile=profile.cov github.com/goci-io/deployment-webhook/cmd/...
	# @TODO implement multiple coverage profiles

run/local:
	./bin/webhook-server
