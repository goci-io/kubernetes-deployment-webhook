ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
.DEFAULT_GOAL := bin/in-docker

export GO_VERSION ?= 1.14
export GO_APP = github.com/goci-io/deployment-webhook

export CONFIG_DIR ?= ./config
export WEBHOOK_SECRET ?= test
export FORCE_NON_TLS_SERVER ?= 1
export ORGANIZATION_WHITELIST ?= goci-io

bin/in-docker:
	docker run --rm \
		-v $(ROOT_DIR):/usr/src/$(GO_APP) \
		golang:$(GO_VERSION) \
		make -C /usr/src/$(GO_APP) bin/server

bin/server:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

bin/server/darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

run:
	go build

tests:
	go test $(GO_APP)/cmd/server/...
	go test $(GO_APP)/cmd/kubernetes/...

coverage:
	go test -v -coverprofile=profile.cov $(GO_APP)/cmd/...
	# @TODO implement multiple coverage profiles

run/local:
	./bin/webhook-server
