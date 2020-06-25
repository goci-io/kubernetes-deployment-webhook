.DEFAULT_GOAL := image/helm-deployment

image/helm-deployment:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./webhook-server ./cmd/helm

image/helm-deployment/darwin:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags="-s -w" -o ./webhook-server ./cmd/helm

run:
	go build

tests:
	go test github.com/goci-io/deployment-webhook/cmd/helm/...
