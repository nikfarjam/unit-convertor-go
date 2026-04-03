# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.26.1-alpine3.23 AS build-stage

WORKDIR /app

COPY go.mod *.go ./
COPY ./pkg ./pkg
COPY ./version ./

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:3.23.3 AS run-stage

WORKDIR /app

RUN apk update --no-cache && apk upgrade && \
    apk cache clean && \
    rm -rf /var/cache/apk/* && \
    adduser -D -u 1001 app && \
    mkdir -p /var/log/unit-converter && chown app:app /var/log/unit-converter
USER app:app

COPY --from=build-stage /app/api-server /app/version /app/

ENV LOG_LEVEL=INFO \
    LOG_OUTPUT=STANDARD \
    LOG_FILE_PATH=/var/log/unit-converter/app.log
EXPOSE 9090

ENTRYPOINT ["/app/api-server"]