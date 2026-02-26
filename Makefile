init:
	go mod tidy

build:
	go build -o ./bin/

run:
	go run main.go

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

docker-image:
	# pass ARCH if you need a specific architecture (e.g. ARCH=arm64/)
	# defaults to empty which results in standard images (amd64 on most hosts)
	docker build --build-arg ARCH=$(ARCH) -t unit-converter-api:$(shell date -I) -f Dockerfile .
	@echo "\nDocker image 'unit-converter-api:$(shell date -I)' built successfully."

clean:
	rm -rf ./bin
	rm -f coverage.out

.PHONY: init build run test coverage lint docker-image clean