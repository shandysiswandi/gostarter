# ***** *****
# VARIABLES
# ***** *****
# DRIVER: Specifies the database driver. Defaults to 'mysql'; can be overridden by setting DRIVER=postgres in the command.
# DB_CONN: Constructs the database connection string based on the DRIVER type (MySQL or PostgreSQL).
# DB_USR and DB_PWD: Set as environment variables for secure database authentication.
DRIVER ?= mysql
DB_CONN ?= $(DB_USR):$(DB_PWD)@tcp(localhost:3306)/gostarter?parseTime=true

# If DRIVER is set to 'postgres', updates DB_CONN for PostgreSQL compatibility.
ifeq ($(DRIVER), postgres)
	DB_CONN = "postgres://${DB_USR}:${DB_PWD}@localhost:5432/gostarter?sslmode=disable"
endif

# ***** *****
# DEVELOPMENT
# ***** *****
# install: Installs CLI tools necessary for development (mockery, buf, goose, reflex, golangci-lint).
install:
	@go install github.com/vektra/mockery/v2@v2.46.3
	@go install github.com/bufbuild/buf/cmd/buf@v1.46.0
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/cespare/reflex@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

# run: Watches for changes in '.go' files, recompiling and rerunning the main application with reflex for hot reloading.
run:
	@reflex -r '\.go$$' -s -- go run main.go

# lint: Runs golangci-lint for static analysis, ensuring Go code quality and adherence to best practices.
lint:
	@golangci-lint run

# load: Sends concurrent gRPC requests for load testing
load:
	@k6 run script/k6/index.js

# Code Quality
scan: test-unit
	@sonar-scanner


# ***** *****
# TESTING
# ***** *****
# test-unit: Runs unit tests excluding directories like mocker, mockz, and app. Produces a coverage report and detects concurrency issues.
test-unit:
	@go test $(shell go list ./pkg/... ./internal/... | grep -vE '/mocker|/mockz|/pkg/goerror/pb|/app') \
	-coverprofile=coverage.out -parallel 4 -race
	@go tool cover -func=coverage.out | grep total
	@if [ "$(HTML)" = "true" ]; then \
		go tool cover -html=coverage.out -o=cover.html; \
		open cover.html; \
	fi

# test-integration: Runs integration tests located in the './tests/...' directory using up to 4 parallel tests.
test-integration:
	@go test ./tests/... -parallel 4

# ***** *****
# GO TOOL
# ***** *****
# tidy: Cleans up Go module dependencies by removing any unused packages.
tidy:
	@go mod tidy

# update: Updates all project dependencies to their latest versions.
update:
	@go get -u ./...

# ***** *****
# GENERATOR
# ***** *****
# gen-proto: Generates Protobuf files from the 'api' directory using buf, overwriting any previous 'gen-proto' directory.
gen-proto:
	@cd api && rm -rf gen-proto && buf dep update && buf build && buf generate && cd ..

# gen-gql: Generates GraphQL files in the 'api' directory using gqlgen, clearing any existing 'gen-gql' directory.
gen-gql:
	@cd api && rm -rf gen-gql && go run github.com/99designs/gqlgen@v0.17.49 generate

# gen-mock: Generates mock files for interfaces using mockery, facilitating mock-based testing.
gen-mock:
	@mockery

# gen-rsa: Creates an RSA key pair for encryption. Outputs keys as both PEM files and Base64-encoded files in the 'config' directory.
gen-rsa:
	@cd config && \
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048 && \
	openssl rsa -pubout -in private_key.pem -out public_key.pem && \
	openssl base64 -A -in private_key.pem -out private_key_base64.txt && \
	openssl base64 -A -in public_key.pem -out public_key_base64.txt

# ***** *****
# MIGRATION
# ***** *****
# migration-create: Creates a new SQL migration file in the 'migration' directory using goose.
migration-create:
	@goose -dir migrations create example sql

# migration-up: Applies all pending migrations for the specified database driver, MySQL or PostgreSQL.
migration-up:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" up

# migration-down: Rolls back the last applied migration for the specified database driver.
migration-down:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" down

# migration-reset: Rolls back all migrations and reapplies them for the specified driver to reset the database state.
migration-reset:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" reset

# ***** *****
# DOCKER
# ***** *****
# docker-build: Builds a Docker image for the application with a timezone argument, tagging the image as 'gostarter'.
docker-build:
	@docker build --build-arg TZ="Asia/Jakarta" -t gostarter .

# docker-run: Runs the Docker container named 'gostarter-container', mapping ports 8081 and 50001, and mounts the config file.
docker-run:
	@docker run --name gostarter-container \
	-p 8081:8081 -p 50001:50001 \
	-v ./config/config.docker.yaml:/config/config.yaml \
	gostarter

# docker-rm: Removes the Docker container named 'gostarter-container'.
docker-rm:
	@docker rm gostarter-container
