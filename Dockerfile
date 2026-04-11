# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.26.2-alpine3.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o api-server main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/static-debian12 AS production

WORKDIR /app

COPY --from=build-stage /app/api-server /app/version /app/

ENV LOG_LEVEL=INFO \
    LOG_OUTPUT=STANDARD \
    LOG_FILE_PATH=/tmp/unit-converter.log
EXPOSE 9090

ENTRYPOINT ["/app/api-server"]