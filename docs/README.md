# Documentation

This documentation includes design and user information in addition to your GoDoc-generated documentation.

## CLI Tools

The following command-line tools are used:

- [Mockery](https://github.com/vektra/mockery)
- [GolangCI-Lint](https://github.com/golangci/golangci-lint)
- [Goose](https://github.com/pressly/goose)
- [Reflex](https://github.com/cespare/reflex)
- [Buf](https://github.com/bufbuild/buf)

## Running the Application

1. **Configuration**: Create a file named `config/config.yaml`. For reference, use the example file `config/config.example.yaml`.

2. **Install CLI Tools**: Download and install the CLI tools listed above.

3. **Database Setup**: Set up the MySQL database as needed.

4. **Migrate Database**: Run the CLI command `make migrate-up`. Ensure that the environment variables `$DB_USR` and `$DB_PWD` are set.

5. **Run the Server**: Start the server with the CLI command `make run`.

### Optional Development Commands

- **Generate Mocks**: Run `make mock` to generate mocks. Ensure you have configured `.mockery.yml`.

- **Generate Protobuf**: Run `make proto` to generate protobuf.

- **Generate GraphQL**: Run `make gql` to generate graphQL.

- **Run Linter**: Check code quality with `make lint`.

## Libraries Used

- [gen-go-protovalidate](https://buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go)
- [gogle-pubsub](https://cloud.google.com/go/pubsub)
- [gqlgen](https://github.com/99designs/gqlgen)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [protovalidate-go](https://github.com/bufbuild/protovalidate-go)
- [snowflake](https://github.com/bwmarrin/snowflake)
- [validator](https://github.com/go-playground/validator/v10)
- [mysql](https://github.com/go-sql-driver/mysql)
- [uuid](https://github.com/google/uuid)
- [httprouter](https://github.com/julienschmidt/httprouter)
- [koanf-parser-yaml](https://github.com/knadh/koanf/parsers/yaml)
- [koanf-provider-file](https://github.com/knadh/koanf/providers/file)
- [koanf](https://github.com/knadh/koanf/v2)
- [go-nanoid](https://github.com/matoous/go-nanoid/v2)
- [go-redis](https://github.com/redis/go-redis/v9)
- [cors](https://github.com/rs/cors)
- [logrus](https://github.com/sirupsen/logrus)
- [viper](https://github.com/spf13/viper)
- [testify](https://github.com/stretchr/testify)
- [gqlparser](https://github.com/vektah/gqlparser/v2)
- [msgpack](https://github.com/vmihailenco/msgpack/v5)
- [zap](https://go.uber.org/zap)
- [google-api](https://google.golang.org/api)
- [genproto](https://google.golang.org/genproto)
- [grpc](https://google.golang.org/grpc)
- [protobuf](https://google.golang.org/protobuf)
