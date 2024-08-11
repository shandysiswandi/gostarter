run:
	@reflex -r '\.go$$' -s -- go run main.go

lint:
	@golangci-lint run

test:
	@go test ./... -parallel 4 -cover -v

mock:
	@mockery

migration-create:
	@goose -dir migration create example sql

migration-up:
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" up

migration-down:
	@goose -dir migration mysql "${DB_USR}:${DB_PWD}@tcp(localhost:3306)/gostarter?parseTime=true" down