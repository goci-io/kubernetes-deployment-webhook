.DEFAULT_GOAL := image/server

export RELEASE ?= latest
export CONFIG_DIR ?= ./config
export WEBHOOK_SECRET ?= test
export FORCE_NON_TLS_SERVER ?= 1
export ORGANIZATION_WHITELIST ?= goci-io

image/server:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

image/server/darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w" -o ./bin/webhook-server ./cmd/server

image/docker: image/server
	docker build kubernetes-deployment-webhook .

image/docker/release:
	docker tag kubernetes-deployment-webhook docker.pkg.github.com/goci-io/kubernetes-deployment-webhook/server:$(RELEASE)
	docker push docker.pkg.github.com/goci-io/kubernetes-deployment-webhook/server:$(RELEASE)

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
