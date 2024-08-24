run:
	@reflex -r '\.go$$' -s -- go run main.go

lint:
	@golangci-lint run

test:
	@go test ./... -coverprofile=coverage.out -parallel 4
	@go tool cover -func=coverage.out | grep total

mock:
	@mockery

tidy:
	@go mod tidy

proto:
	@cd api && rm -rf gen-proto && buf mod update && buf build && buf generate && cd ..

gql:
	@cd api && rm -rf gen-gql && go run github.com/99designs/gqlgen@v0.17.48 generate

check: proto gql mock tidy lint test run

migration-create:
	@goose -dir migration create example sql

migration-up:
	@goose -dir migration fix
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" up

migration-down:
	@goose -dir migration fix
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" down

docker-build:
	@docker build --build-arg TZ="Asia/Jakarta" -t gostarter .

docker-run:
	@docker run --name gostarter-container \
	-p 8081:8081 -p 50001:50001 \
	-v ./config/config.docker.yaml:/config/config.yaml \
	gostarter

docker-rm:
	@docker rm gostarter-container