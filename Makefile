# ***** *****
# VARIABLES
# ***** *****
# DRIVER: Specifies the database driver. Defaults to 'mysql', you can override this by passing DRIVER=postgres in the command.
# DB_CONN: The database connection string, constructed based on the DRIVER type (MySQL or PostgreSQL).
# DB_USR and DB_PWD are expected to be set as environment variables for database authentication.
DRIVER ?= mysql
DB_CONN ?= $(DB_USR):$(DB_PWD)@tcp(localhost:3306)/gostarter?parseTime=true

# If the database driver is 'postgres', update the connection string accordingly.
ifeq ($(DRIVER), postgres)
	DB_CONN = "postgres://${DB_USR}:${DB_PWD}@localhost:5432/gostarter?sslmode=disable"
endif

# ***** *****
# DEVELOPMENT 
# ***** *****
#
# 'run': Watches for file changes with the '.go' extension and automatically recompiles and runs the main Go application. It uses the 'reflex' tool for hot reloading.
run:
	@reflex -r '\.go$$' -s -- go run main.go

# 'lint': Runs 'golangci-lint' to perform static code analysis, ensuring code quality and adherence to Go standards.
lint:
	@golangci-lint run

# ***** *****
# TESTING 
# ***** *****
#
# 'test-unit': Runs unit tests, excluding certain directories (e.g., mocker, mockz, app). It produces a coverage report, uses parallelism (up to 4 tests), and includes race detection for concurrency issues.
test-unit:
	@go test $(shell go list ./pkg/... ./internal/... | grep -vE '/mocker|/mockz|/app') \
	-coverprofile=coverage.out -parallel 4 -race
	@go tool cover -func=coverage.out | grep total

# 'test-integration': Runs integration tests located in the './tests/...' directory using 4 parallel test instances.
test-integration:
	@go test ./tests/... -parallel 4

# ***** *****
# GO TOOL 
# ***** *****
#
# 'tidy': Cleans up the Go module by removing unused dependencies with 'go mod tidy'.
tidy:
	@go mod tidy

# 'update': Updates all the Go dependencies in the project to their latest versions.
update:
	@go get -u ./...

# ***** *****
# GENERATOR
# ***** *****
#
# 'gen-proto': Generates Protobuf files from the 'api' directory using the 'buf' tool. Removes the previous 'gen-proto' directory before regenerating.
gen-proto:
	@cd api && rm -rf gen-proto && buf dep update && buf build && buf generate && cd ..

# 'gen-gql': Generates GraphQL files using the 'gqlgen' tool in the 'api' directory. Removes the previous 'gen-proto' directory before regenerating.
gen-gql:
	@cd api && rm -rf gen-gql && go run github.com/99designs/gqlgen@v0.17.49 generate

# 'gen-mock': Generates mock files using 'mockery' for mocking interfaces in tests.
gen-mock:
	@mockery

# 'gen-rsa': Generates an RSA key pair for encryption purposes in the 'config' directory. Outputs the private and public keys both as PEM files and in Base64 format.
gen-rsa:
	@cd config && \
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048 && \
	openssl rsa -pubout -in private_key.pem -out public_key.pem && \
	openssl base64 -A -in private_key.pem -out private_key_base64.txt && \
	openssl base64 -A -in public_key.pem -out public_key_base64.txt

# ***** *****
# MIGRATION 
# ***** *****
#
# 'migration-create': Creates a new database migration file using 'goose' in the 'migration' directory.
migration-create:
	@goose -dir migration create example sql

# 'migration-up': Applies the database migrations (up direction) for the specified driver (MySQL or PostgreSQL). The migration scripts are located in the 'migrations/<DRIVER>' folder.
migration-up:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" up

# 'migration-down': Rolls back the latest database migration (down direction) for the specified driver.
migration-down:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" down

# 'migration-reset': Resets the database by rolling back all migrations and then reapplying them for the specified driver.
migration-reset:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" reset

# ***** *****
# DOCKER 
# ***** *****
#
# 'docker-build': Builds a Docker image for the application using the specified timezone (Asia/Jakarta) as an argument. The resulting image is tagged as 'gostarter'.
docker-build:
	@docker build --build-arg TZ="Asia/Jakarta" -t gostarter .

# 'docker-run': Runs the application in a Docker container with the name 'gostarter-container'. It maps ports 8081 and 50001 from the container to the host and mounts a configuration file from the host system.
docker-run:
	@docker run --name gostarter-container \
	-p 8081:8081 -p 50001:50001 \
	-v ./config/config.docker.yaml:/config/config.yaml \
	gostarter

# 'docker-rm': Removes the Docker container named 'gostarter-container'.
docker-rm:
	@docker rm gostarter-container