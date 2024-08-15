run:
	@reflex -r '\.go$$' -s -- go run main.go

lint:
	@golangci-lint run

test:
	@go test ./... -coverprofile=coverage.out -parallel 4

mock:
	@mockery
	@rm pkg/pkgmock/mock_client_option.go

tidy:
	@go mod tidy

proto:
	@cd api && buf mod update && buf build && buf generate && cd ..

check: proto mock tidy lint test

migration-create:
	@goose -dir migration create example sql

migration-up:
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" up

migration-down:
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" down