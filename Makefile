IMAGE_TAG ?= unit-converter-api-lightweight:$(shell date -I)

init:
	go mod tidy

build:
	go build -o ./bin/

run:
	go run main.go

run-dev:
	LOG_LEVEL=DEBUG go run main.go

test:
	go test -v ./...

coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

coverage-ci:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'

coverage-codecov:
	go test -coverprofile=coverage.txt ./...

lint:
	go mod verify
	go fmt ./...

BUILDX_ARGS ?= --load --attest type=sbom --attest type=provenance
build-image:
	@echo "\nBuild Docker image '$(IMAGE_TAG)' with SBOM."
	docker buildx build $(BUILDX_ARGS) -t $(IMAGE_TAG) -f Dockerfile .

clean:
	rm -rf ./bin
	rm -f coverage.out

.PHONY: init build run run-dev test coverage-html coverage-ci coverage-codecov lint build-image clean