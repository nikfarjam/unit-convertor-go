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

LABEL "com.datadoghq.ad.check_names"='["unit_converter"]'
LABEL "com.datadoghq.ad.init_configs"='[{}]'
LABEL "com.datadoghq.ad.instances"='[{"unit_converter": "http://%%host%%:%%port%%/version"}]'
LABEL "com.datadoghq.ad.logs"='[{"source": "unit_converter", "service": "api-server", "type": "http"}]'

ENTRYPOINT ["/app/api-server"]