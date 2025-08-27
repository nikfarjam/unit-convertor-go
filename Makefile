init:
	go mod tidy
	go mod verify
	go fmt .

build:
	go build -o ./bin/

run:
	go run main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf ./bin
	rm -f coverage.out

.PHONY: init build run test coverage lint