# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.25.0 AS build-stage

WORKDIR /app

COPY go.mod *.go ./
COPY ./pkg ./pkg

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:3.22 AS build-release-stage

WORKDIR /

RUN adduser -D -u 1001 nonroot
USER nonroot:nonroot

COPY --from=build-stage /app/api-server /api-server

EXPOSE 9090

ENTRYPOINT ["/api-server"]