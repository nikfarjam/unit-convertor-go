# Unit Convertor

An example serverless RESTful API written in GoLang

## Overview

This project is a simple unit conversion service that provides an API for converting between temperature units (Fahrenheit and Celsius). It's designed to be lightweight and easy to deploy as a serverless application.

## Features

- Converts between Fahrenheit and Celsius units.
- RESTful API endpoints.
- Written in GoLang for performance and simplicity.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [curl](https://curl.se/download.html)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/unit-convertor-go.git
   cd unit-convertor-go
   ```

2. Project setup and usage via Makefile:

| Command | Description |
| -------------------- | ----------------------------------------- |
| `make init` | Initialize Go modules |
| `make build` | Build the binary in `./bin/` |
| `make run` | Run the application |
| `make test` | Run all tests |
| `make coverage-html` | Generate HTML coverage report |
| `make coverage-ci` | Output total coverage (CI) |
| `make coverage-codecov` | Generate Codecov coverage file |
| `make lint` | Verify modules and format code |
| `make docker-image` | Build Docker image (you can set `ARCH` to your machine architecture; e.g. `make docker-image ARCH=arm64/`) |
| `make clean` | Clean build and coverage files |

1. Run application:

   ```bash
   make build
   ./bin/unit-convertor-go
   ```

### API Usage

The API exposes a single endpoint:

- **POST** `/converter`
   - Request body (JSON):
      - `value`: number (temperature value)
      - `from`: "CELSIUS" or "FAHRENHEIT" (case-insensitive)
      - `to`: "CELSIUS" or "FAHRENHEIT" (case-insensitive)

#### Example curl commands

**Convert Celsius to Fahrenheit:**

```bash
curl -d '{ "value": 25, "from": "CELSIUS", "to": "FAHRENHEIT" }' http://localhost:9090/converter
```

**Convert Fahrenheit to Celsius:**

```bash
curl -d '{ "value": 77, "from": "FAHRENHEIT", "to": "CELSIUS" }' http://localhost:9090/converter
```

**Invalid Unit (returns error):**

```bash
curl -v -d '{ "value": 20, "from": "Kelvin", "to": "Celsius" }' http://localhost:9090/converter
```

#### Error Handling

- If `from` or `to` is not "CELSIUS" or "FAHRENHEIT", the API returns HTTP 400 with an error message.
- If the request body is malformed, the API returns HTTP 400.

### Using the client.sh Script

The `client.sh` script is a simple CLI tool for calling the API.

**Usage:**

```bash
chmod +x client.sh
./client.sh [F|C] [degree_value]
```

- `F`: Convert Fahrenheit to Celsius (input is Fahrenheit)
- `C`: Convert Celsius to Fahrenheit (input is Celsius)
- `degree_value`: Numeric temperature value to convert

**Examples:**

```bash
./client.sh F 77
./client.sh C 25
```

If you provide an invalid unit or a non-numeric value, the script will print an error and exit.

## License

This project is licensed under the Apache-2.0 License.
