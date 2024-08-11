# Documentation

This documentation includes design and user information in addition to your GoDoc-generated documentation.

## CLI Tools

The following command-line tools are used:

- [Mockery](https://github.com/vektra/mockery)
- [GolangCI-Lint](https://github.com/golangci/golangci-lint)
- [Goose](https://github.com/pressly/goose)
- [Reflex](https://github.com/cespare/reflex)

## Running the Application

1. **Configuration**: Create a file named `config/config.yaml`. For reference, use the example file `config/config.example.yaml`.

2. **Install CLI Tools**: Download and install the CLI tools listed above.

3. **Database Setup**: Set up the MySQL database as needed.

4. **Migrate Database**: Run the CLI command `make migrate-up`. Ensure that the environment variables `$DB_USR` and `$DB_PWD` are set.

5. **Run the Server**: Start the server with the CLI command `make run`.

### Optional Development Commands

- **Generate Mocks**: Run `make mock` to generate mocks. Ensure you have configured `.mockery.yml`.

- **Run Linter**: Check code quality with `make lint`.

## Libraries Used

- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [snowflake](https://github.com/bwmarrin/snowflake)
- [validator](https://github.com/go-playground/validator/v10)
- [mysql](https://github.com/go-sql-driver/mysql)
- [uuid](https://github.com/google/uuid)
- [httprouter](https://github.com/julienschmidt/httprouter)
- [go-nanoid](https://github.com/matoous/go-nanoid/v2)
- [go-redis](https://github.com/redis/go-redis/v9)
- [cors](https://github.com/rs/cors)
- [viper](https://github.com/spf13/viper)
- [testify](https://github.com/stretchr/testify)
- [msgpack](https://github.com/vmihailenco/msgpack/v5)
