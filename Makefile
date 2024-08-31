DRIVER ?= mysql
DB_CONN ?= $(DB_USR):$(DB_PWD)@tcp(localhost:3306)/gostarter?parseTime=true

ifeq ($(DRIVER), postgres)
	DB_CONN = "postgres://${DB_USR}:${DB_PWD}@localhost:5432/gostarter?sslmode=disable"
endif


run:
	@reflex -r '\.go$$' -s -- go run main.go

lint:
	@golangci-lint run

test:
	@go test $(go list ./pkg/... ./internal/... | grep -vE "/mocker$|/mockz$") -coverprofile=coverage.out -parallel 4
	@go tool cover -func=coverage.out | grep total

test-integration:
	@go test ./tests/... -parallel 4

mock:
	@mockery

tidy:
	@go mod tidy

update:
	@go get -u ./...

proto:
	@cd api && rm -rf gen-proto && buf dep update && buf build && buf generate && cd ..

gql:
	@cd api && rm -rf gen-gql && go run github.com/99designs/gqlgen@v0.17.49 generate

check: proto gql mock tidy lint test run

migration-create:
	@goose -dir migration create example sql

migration-up:
	@goose -dir migration/$(DRIVER) fix
	@goose -dir migration/$(DRIVER) $(DRIVER) "$(DB_CONN)" up

migration-down:
	@goose -dir migration/$(DRIVER) fix
	@goose -dir migration/$(DRIVER) $(DRIVER) "$(DB_CONN)" down

docker-build:
	@docker build --build-arg TZ="Asia/Jakarta" -t gostarter .

docker-run:
	@docker run --name gostarter-container \
	-p 8081:8081 -p 50001:50001 \
	-v ./config/config.docker.yaml:/config/config.yaml \
	gostarter

docker-rm:
	@docker rm gostarter-container