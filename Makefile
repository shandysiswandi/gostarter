DRIVER ?= mysql
DB_CONN ?= $(DB_USR):$(DB_PWD)@tcp(localhost:3306)/gostarter?parseTime=true

ifeq ($(DRIVER), postgres)
	DB_CONN = "postgres://${DB_USR}:${DB_PWD}@localhost:5432/gostarter?sslmode=disable"
endif

install:
	@go install github.com/vektra/mockery/v2@v2.50.0
	@go install github.com/bufbuild/buf/cmd/buf@v1.47.2
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/cespare/reflex@latest
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

tools:
	@go install golang.org/x/tools/gopls@latest
	@go install github.com/cweill/gotests/gotests@latest
	@go install github.com/fatih/gomodifytags@latest
	@go install github.com/josharian/impl@latest
	@go install github.com/haya14busa/goplay/cmd/goplay@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

run:
	@reflex -r '\.go$$' -s -- sh -c "LOCAL=true go run main.go"

lint:
	@golangci-lint run

test:
	@go test $(shell go list ./pkg/... ./internal/... | grep -vE '/mocker|/mockz|/pkg/goerror/pb|/app') \
	-coverprofile=coverage.out -race -parallel 4
	@go tool cover -func=coverage.out | grep total
	@if [ "$(HTML)" = "true" ]; then \
		go tool cover -html=coverage.out -o=cover.html; \
		open cover.html; \
	fi

tidy:
	@go mod tidy

update:
	@go get -u ./...

gen-proto:
	@cd api && rm -rf gen-proto && buf dep update && buf build && buf generate && cd ..

gen-gql:
	@cd api && rm -rf gen-gql && go run github.com/99designs/gqlgen@v0.17.60 generate

gen-mock:
	@mockery

gen-rsa:
	@cd config && \
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048 && \
	openssl rsa -pubout -in private_key.pem -out public_key.pem && \
	openssl base64 -A -in private_key.pem -out private_key_base64.txt && \
	openssl base64 -A -in public_key.pem -out public_key_base64.txt

migration-create:
	@goose -dir migrations create example sql

migration-up:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" up

migration-down:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" down

migration-reset:
	@goose -dir migrations/$(DRIVER) fix
	@goose -dir migrations/$(DRIVER) $(DRIVER) "$(DB_CONN)" reset


docker-build:
	@docker build --build-arg TZ="Asia/Jakarta" -t gostarter .

docker-build-goose:
	@docker build --build-arg GOOSE_VERSION_TAG="v3.23.0" -f Dockerfile.goose -t goose:v3.23.0 .

compose-up:
	@docker compose up --build
	
compose-down:
	@docker compose down 