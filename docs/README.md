# Project Documentation

This document provides a detailed guide for setting up, running, and developing this application, covering key configuration and CLI tools required.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Running the Application](#running-the-application)
4. [Development Commands](#development-commands)
5. [Libraries Used](#libraries-used)
6. [Contact](#contact)

---

## Prerequisites

The application requires several command-line tools and dependencies:

- **Mockery** - Used for generating mocks ([GitHub link](https://github.com/vektra/mockery))
- **GolangCI-Lint** - Linting tool for Go ([GitHub link](https://github.com/golangci/golangci-lint))
- **Goose** - Migration tool for database management ([GitHub link](https://github.com/pressly/goose))
- **Reflex** - Automatic reload of Go code on change ([GitHub link](https://github.com/cespare/reflex))
- **Buf** - Protobuf generator and linter ([GitHub link](https://github.com/bufbuild/buf))

Ensure these tools are installed prior to running the application.

## Installation

### 1. Configuration

1. **Configuration File**: Create a configuration file at `config/config.yaml` based on the provided template `config/config.example.yaml`.
2. **Environment Variables**: Set necessary environment variables such as `$DB_USR` and `$DB_PWD` for database connections.

### 2. Database Setup

- **Supported Databases**: The application supports MySQL and PostgreSQL databases.
- **Migration Setup**: Use Goose to run migrations.

### 3. Install CLI Tools

If the CLI tools listed above are not already installed, follow the instructions on their respective GitHub repositories to install them.

## Running the Application

1. **Database Migration**: Migrate the database schema with:

   ```sh
   make migration-up
   ```

   - The command defaults to using the MySQL driver. For PostgreSQL, use:
     ```sh
     make migration-up DRIVER=postgres
     ```
     Ensure `$DB_USR` and `$DB_PWD` environment variables are set.

2. **Start the Server**: Launch the server with:
   ```sh
   make run
   ```
   - The server uses Reflex to automatically reload upon code changes.

## Development Commands

The following commands are available for development:

- **Generate Mocks**:

  ```sh
  make gen-mock
  ```

  Generates mocks using Mockery. Ensure `.mockery.yml` is configured as needed.

- **Generate Protobuf**:

  ```sh
  make gen-proto
  ```

  Generates Protobuf files using Buf.

- **Generate GraphQL**:

  ```sh
  make gen-gql
  ```

  Generates GraphQL schema.

- **Generate RSA Keys**:

  ```sh
  make gen-rsa
  ```

  Generates RSA keys.

- **Run Linter**:
  ```sh
  make lint
  ```
  Runs GolangCI-Lint to check code quality.

## Libraries Used

The following libraries are essential to the application’s functionality:

| Library                                                      | Version  | Purpose                               |
| ------------------------------------------------------------ | -------- | ------------------------------------- |
| `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go` | v1.34.2  | Protocol buffer validation            |
| `cloud.google.com/go/pubsub`                                 | v1.43.0  | Google Cloud Pub/Sub client           |
| `github.com/99designs/gqlgen`                                | v0.17.49 | GraphQL server generator              |
| `github.com/DATA-DOG/go-sqlmock`                             | v1.5.2   | SQL mock for testing                  |
| `github.com/bufbuild/protovalidate-go`                       | v0.6.5   | Protobuf validation library           |
| `github.com/bwmarrin/snowflake`                              | v0.3.0   | Distributed unique ID generator       |
| `github.com/doug-martin/goqu/v9`                             | v9.19.0  | SQL builder for Go                    |
| `github.com/go-playground/validator/v10`                     | v10.22.1 | Input validation library              |
| `github.com/go-redis/redismock/v9`                           | v9.2.0   | Redis mock for testing                |
| `github.com/go-resty/resty/v2`                               | v2.14.0  | HTTP and REST client                  |
| `github.com/golang-jwt/jwt/v4`                               | v4.5.0   | JWT implementation                    |
| `github.com/google/uuid`                                     | v1.6.0   | UUID generator                        |
| `github.com/julienschmidt/httprouter`                        | v1.3.0   | High-performance HTTP request router  |
| `github.com/knadh/koanf`                                     | v2.1.1   | Library for reading configuration     |
| `github.com/lib/pq`                                          | v1.10.9  | PostgreSQL driver for Go              |
| `github.com/matoous/go-nanoid/v2`                            | v2.1.0   | NanoID generator                      |
| `github.com/redis/go-redis/v9`                               | v9.6.1   | Redis client for Go                   |
| `github.com/rs/cors`                                         | v1.11.1  | CORS handling middleware              |
| `github.com/spf13/viper`                                     | v1.19.0  | Go configuration library              |
| `github.com/stretchr/testify`                                | v1.9.0   | Testing utilities                     |
| `github.com/vektah/gqlparser/v2`                             | v2.5.16  | GraphQL parser                        |
| `github.com/vmihailenco/msgpack/v5`                          | v5.4.1   | MessagePack encoding library          |
| `go.opentelemetry.io/otel`                                   | v1.30.0  | OpenTelemetry instrumentation         |
| `go.uber.org/zap`                                            | v1.27.0  | Logging library                       |
| `golang.org/x/crypto`                                        | v0.27.0  | Cryptography libraries                |
| `google.golang.org/api`                                      | v0.197.0 | Google API client                     |
| `google.golang.org/grpc`                                     | v1.66.1  | gRPC client and server implementation |
| `google.golang.org/protobuf`                                 | v1.34.2  | Protobuf support library              |

## Contact

For questions or support, please contact [Project Maintainer](mailto:shandysiswandi@gmail.com).

---
