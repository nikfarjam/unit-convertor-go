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
- Docker

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/unit-convertor-go.git
   cd unit-convertor-go
   ```

2. Initialize project `make init`
3. Test `make test`
4. Build project `make build`
5. Run application `./bin/unit-convertor-go`

### Example Usage

Here are some example curl commands to interact with the API:

**Convert Celsius to Fahrenheit:**

   ```bash
   curl -d "{ \"value\": 25, \"from\": \"CELSIUS\", \"to\": \"FAHRENHEIT\" }" http://localhost:9090/converter
   ```

**Convert Fahrenheit to Celsius:**

   ```bash
   curl -d "{ \"value\": 77, \"from\": \"FAHRENHEIT\", \"to\": \"CELSIUS\" }" http://localhost:9090/converter
   ```

**Invalid Unit:**

  ```bash
  curl -v -d "{ \"value\": 20, \"from\": \"Kelvin\", \"to\": \"Celsius\" }" http://localhost:9090/converter
  ```

### Using the client.sh Script

The `client.sh` script provides a convenient way to interact with the API from the command line.

**Usage:**

  ```bash
  chmod +x client.sh
  [client.sh](http://_vscodecontentref_/0) [F|C] [degree_value]
  ```

- F: Convert Celsius to Fahrenheit.
- C: Convert Fahrenheit to Celsius.
- degree_value: The temperature value to convert.
Examples:

  ```bash
  client.sh F 77
  client.sh C 25
  ```
