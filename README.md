# Unit Convertor

A lightweight RESTful API written in Go for unit conversion

## Overview

This project is a simple temperature conversion service that provides RESTful API endpoints for converting between Fahrenheit and Celsius. It's designed to be lightweight, performant, and easy to deploy using Docker.

## Features

- Converts between Fahrenheit and Celsius units
- RESTful API with JSON request/response
- Version endpoint for API versioning
- Structured logging with configurable levels
- Docker containerization support
- CLI client script for easy testing

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (1.25 or later)
- [curl](https://curl.se/download.html) (for API testing)
- [Docker](https://www.docker.com/) (optional, for containerized deployment)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/nikfarjam/unit-convertor-go.git
   cd unit-convertor-go
   ```

2. Project setup and usage via Makefile:

| Command | Description |
| -------------------- | ----------------------------------------- |
| `make init` | Initialize Go modules |
| `make build` | Build the binary in `./bin/` |
| `make run` | Run the application |
| `make run-dev` | Run the application with debug log |
| `make test` | Run all tests |
| `make coverage-html` | Generate HTML coverage report |
| `make coverage-ci` | Output total coverage (CI/CD) |
| `make coverage-codecov` | Generate Codecov coverage file |
| `make lint` | Verify modules and format code |
| `make build-image` | Build Docker image |
| `make clean` | Clean build and coverage files |

3. Run application:

   ```bash
   make build
   ./bin/unit-convertor-go
   ```

## API Documentation

The API exposes two endpoints:

### POST `/converter`

Converts temperature between Celsius and Fahrenheit.

**Request Body (JSON):**

```json
{
  "value": 25,
  "from": "CELSIUS",
  "to": "FAHRENHEIT"
}
```

**Response (JSON):**

```json
{
  "value": 77,
  "unit": "FAHRENHEIT"
}
```

**Example curl commands:**

Convert Celsius to Fahrenheit:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"value": 25, "from": "CELSIUS", "to": "FAHRENHEIT"}' \
  http://localhost:9090/converter
```

Convert Fahrenheit to Celsius:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"value": 77, "from": "FAHRENHEIT", "to": "CELSIUS"}' \
  http://localhost:9090/converter
```

### GET `/version`

Returns the current version of the API.

**Response:** Plain text version string (e.g., "v1.1.0")

**Example:**

```bash
curl http://localhost:9090/version
```

## Error Handling

- **400 Bad Request**: Invalid unit names (must be "CELSIUS" or "FAHRENHEIT"), malformed JSON, or invalid request structure
- **500 Internal Server Error**: Server-side processing errors

## Configuration

The application supports the following environment variables:

- `LOG_LEVEL`: Set logging level (DEBUG, INFO, WARN, ERROR) - default: INFO
- `LOG_OUTPUT`: Log output destination (STANDARD or FILE) - default: STANDARD
- `LOG_FILE_PATH`: Path for log file when LOG_OUTPUT=FILE - default: app.log
- `UC_VERSION_PATH`: Path to version file - default: `./version`

## Docker Deployment

Build and run with Docker:

```bash
make build-image
docker run -p 9090:9090 unit-converter-api-lightweight:$(date -I)
```

## Docker Hub Repository

The Docker image is available on Docker Hub: [nikfarjam/unit-converter-api-lightweight](https://hub.docker.com/r/nikfarjam/unit-converter-api-lightweight)

You can pull and run the image directly:

```bash
docker pull nikfarjam/unit-converter-api-lightweight:latest
docker run -p 9090:9090 nikfarjam/unit-converter-api-lightweight:latest
```

## Using the Client Script

The `client.sh` script provides a simple CLI interface for testing the API.

**Usage:**

```bash
chmod +x client.sh
./client.sh [F|C] [degree_value]
```

- `F`: Convert from Fahrenheit to Celsius
- `C`: Convert from Celsius to Fahrenheit
- `degree_value`: Numeric temperature value

**Examples:**

```bash
./client.sh F 77    # Converts 77°F to Celsius
./client.sh C 25    # Converts 25°C to Fahrenheit
```

## Development

### Project Structure

```SHELL
.
├── main.go                # Application entry point
├── main_test.go           # Main package tests
├── client.sh              # CLI client script
├── Dockerfile             # Docker build configuration
├── Makefile               # Build and development tasks
├── go.mod                 # Go module definition
├── version                # Version file
├── pkg/
│   ├── api/               # HTTP handlers
│   │   ├── converter.go
│   │   ├── converter_test.go
│   │   ├── version.go
│   │   └── version_test.go
│   └── converter/         # Business logic
│       ├── converter.go
│       └── converter_test.go
└── bin/                   # Build output directory
```

### Running Tests

```bash
make test                   # Run all tests
make coverage-html          # Generate HTML coverage report
make coverage-ci            # Get total coverage percentage
```

### Code Quality

```bash
make lint                   # Format code and verify modules
```

## License

This project is licensed under the Apache-2.0 License.
