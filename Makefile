.DEFAULT_GOAL := image/server

export CONFIG_DIR ?= ./config
export WEBHOOK_SECRET ?= test
export FORCE_NON_TLS_SERVER ?= 1
export ORGANIZATION_WHITELIST ?= goci-io

image/server:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./webhook-server ./cmd/server

image/server/darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w" -o ./webhook-server ./cmd/server

run:
	go build

tests:
	go test github.com/goci-io/deployment-webhook/cmd/server/...

coverage:
	go test -v -coverprofile=profile.cov github.com/goci-io/deployment-webhook/cmd/server/...

run/local:
	./webhook-server
