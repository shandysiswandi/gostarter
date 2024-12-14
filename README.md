# Go Starter

Go Starter is an opinionated, production-ready server API framework written in Go. It is designed to help developers quickly set up robust, scalable, and secure web services with minimal configuration. With a focus on simplicity, extensibility, and best practices, Go Starter provides a strong foundation for building and deploying production-grade applications.

This framework is suitable for a wide range of use cases, from microservices to monolithic applications. It comes with a set of pre-configured tools, libraries, and patterns to streamline development while ensuring maintainability, scalability, and performance.

## Table of Contents

1. [Application](#application)
2. [Makefile](#makefile)
3. [Contact](#contact)

## Application

### Modules

- Authentication
- Users
- Todos

### Transport

- HTTP
- GraphQL
- gRPC

### Database

- MySQL
- PostgreSQL

### Cache

- Redis

## Makefile

Makefile is designed to streamline development, testing, and deployment processes for the application. It includes commands for managing dependencies, running the application, testing, generating files, database migrations, and Docker operations. Below is a breakdown of the available targets and their purposes.

### Variables

- Specifies the database driver to use.
- Default: `mysql`
- Can be overridden by setting `DRIVER=postgres` in the command.

#### DB_CONN

- Constructs the database connection string based on the `DRIVER` type:
  - MySQL: Default connection string for MySQL.
  - PostgreSQL: Updated dynamically if `DRIVER` is set to `postgres`.

#### DB_USR and DB_PWD

- Environment variables used for database authentication.
- Must be set externally for secure connection handling.

### Targets

#### Development

`install`

Installs CLI tools required for development:

- `mockery`: Mock generation for interfaces.

- `buf`: Protocol Buffers tool.

- `goose`: Database migration tool.

- `reflex`: File watcher for hot reloading.

- `golangci-lint`: Go linter.

`run`

- Watches for changes in `.go` files and re-compiles/reruns the main application using `reflex`.

`lint`

- Runs `golangci-lint` for static code analysis and ensures adherence to Go best practices.

`test`

- Runs unit tests, excluding specific directories (`mocker`, `mockz`, etc.).
- Generates a coverage report and detects concurrency issues.
- If `HTML=true` is set, an HTML coverage report is generated and opened.

#### Go Tooling

`tidy`

- Cleans up Go module dependencies by removing unused packages.

`update`

- Updates all project dependencies to their latest versions.

#### Generators

`gen-proto`

- Generates Protocol Buffer files from the `api` directory using `buf`.
- Clears any existing `gen-proto` directory before generating.

`gen-gql`

- Generates GraphQL files in the `api` directory using `gqlgen`.
- Removes the existing `gen-gql` directory prior to generation.

`gen-mock`

- Generates mock files for interfaces using `mockery`.

`gen-rsa`

- Generates an RSA key pair for encryption.

- Outputs keys in both PEM and Base64 formats to the `config` directory.

#### Database Migrations

`migration-create`

- Creates a new SQL migration file in the `migrations` directory using `goose`.

`migration-up`

- Applies all pending migrations for the specified database driver (`mysql` or `postgres`).

`migration-down`

- Rolls back the last applied migration for the specified database driver.

`migration-reset`

- Rolls back all migrations and reapplies them to reset the database state.

#### Docker

`docker-build`

- Builds a Docker image for the application with a timezone argument (`Asia/Jakarta`).
- Tags the image as `gostarter`.

`docker-build-goose`

- Builds a Docker image for the `goose` migration tool using the specified version tag.

`compose-up`

- Brings up the application stack using Docker Compose.
- Includes building images if necessary.

`compose-down`

- Tears down the Docker Compose stack.

### Usage Examples

Here is an example of how to use the Makefile to set up, develop, test, and deploy the application:

- **Install dependencies**:

  ```sh
  make install
  ```

- **Set up the database credentials**: Export the required environment variables:

  ```sh
  export DB_USR=your_database_user
  export DB_PWD=your_database_password
  ```

- **Apply database migrations**:

  ```sh
  make migration-up
  ```

- **Run the application in development mode with hot reloading**:

  ```sh
  make run
  ```

- **Run unit tests and check code coverage**:

  ```sh
  make test
  ```

- **Perform static code analysis**:

  ```sh
  make lint
  ```

- **Generate required files (e.g., Protobuf, GraphQL, mocks)**:

  ```sh
  make gen-proto
  make gen-gql
  make gen-mock
  ```

- **Build the Docker image**:

  ```sh
  make docker-build
  ```

- **Bring up the application stack using Docker Compose**:

  ```sh
  make compose-up
  ```

- **Tear down the Docker Compose stack**:

  ```sh
  make compose-down
  ```

## Contact

For questions or support, please contact [Project Maintainer](mailto:shandysiswandi@gmail.com).
