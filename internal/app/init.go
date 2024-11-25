// Package app provides a structured framework for initializing, configuring,
// and managing the core components of the application. It handles the setup
// of essential services such as logging, configuration management, database
// connections, Redis client, HTTP router and server, utility libraries, and
// background tasks. The package ensures that all components are properly
// initialized at startup and gracefully shut down when the application stops,
// offering a robust foundation for building and maintaining the application.
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goroutine"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/messaging/googlepubsub"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/instrument"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// initConfig initializes the application's configuration by loading settings from a specified
// YAML file using the Viper library. If the configuration cannot be loaded, the application
// will log a fatal error and terminate. This method should be called early in the application's
// initialization process.
func (a *App) initConfig() {
	cfg, err := config.NewViperConfig("config/config.yaml")
	if err != nil {
		log.Fatalln("failed to init config", err)
	}

	os.Setenv("TZ", cfg.GetString(`tz`))

	a.config = cfg
}

// initTelemetry sets up telemetry for the application, configuring it to use a Zap logger
// with the specified logging level. This enables logging and monitoring capabilities
// across the application, allowing tracking of application metrics and logs for observability.
func (a *App) initTelemetry() {
	filterKeys := []string{"authorization", "password", "access_token", "refresh_token"}

	a.telemetry = telemetry.NewTelemetry(
		telemetry.WithServiceName(a.config.GetString("telemetry.name")),
		telemetry.WithLogFilter(filterKeys...),
		telemetry.WithZapLogger(logger.InfoLevel, filterKeys...),
		telemetry.WithOTLP(a.config.GetString("telemetry.otlp.grpc.address")),
	)
}

// initLibraries initializes various utility libraries used throughout the application,
// such as UID generators, clock, codecs for JSON and MsgPack, and the validation library.
// If any library fails to initialize, the application will log a fatal error and terminate.
func (a *App) initLibraries() {
	snow, err := uid.NewSnowflakeNumber()
	if err != nil {
		log.Fatalln("failed to init uid number snowflake", err)
	}

	pvalidator, err := validation.NewProtoValidator()
	if err != nil {
		log.Fatalln("failed to init validation proto validator", err)
	}

	jewete, err := jwt.NewJSONWebToken(
		a.config.GetString("jwt.private.key"),
		a.config.GetString("jwt.public.key"),
	)
	if err != nil {
		log.Fatalln("failed to init json web token (jwt)", err)
	}

	a.jwt = jewete
	a.uidNumber = snow
	a.protoValidator = pvalidator

	a.clock = clock.New()
	a.uuid = uid.NewUUIDString()
	a.hash = hash.NewBcryptHash(10)
	a.secHash = hash.NewHMACSHA256Hash(a.config.GetString("jwt.hash.secret"))
	a.codecJSON = codec.NewJSONCodec()
	a.goroutine = goroutine.NewManager(100)
	a.codecMsgPack = codec.NewMsgPackCodec()
	a.validator = validation.NewV10Validator()
}

// dsnMySQL constructs a Data Source Name (DSN) for connecting to a MySQL database
// using the application's configuration. It includes connection options such as
// time zone and parseTime settings.
func (a *App) dsnMySQL() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		a.config.GetString(`database.user`),
		a.config.GetString(`database.pass`),
		a.config.GetString(`database.host`),
		a.config.GetString(`database.port`),
		a.config.GetString(`database.name`),
	)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", a.config.GetString(`tz`))
	val.Encode()

	return fmt.Sprintf("%s?%s", dsn, val.Encode())
}

// dsnPostgreSQL constructs a Data Source Name (DSN) for connecting to a PostgreSQL database
// using the application's configuration. It includes connection options such as SSL mode settings.
func (a *App) dsnPostgreSQL() string {
	dsn := fmt.Sprintf("%s:%s@%s:%s/%s",
		a.config.GetString(`database.user`),
		a.config.GetString(`database.pass`),
		a.config.GetString(`database.host`),
		a.config.GetString(`database.port`),
		a.config.GetString(`database.name`),
	)
	val := url.Values{}
	val.Add("sslmode", "disable")
	val.Encode()

	return fmt.Sprintf("postgres://%s?%s", dsn, val.Encode())
}

// initDatabase initializes the application's database connection using settings from the
// configuration. It sets up connection pooling parameters and tests the connection by pinging
// the database. If any step fails, the application will log a fatal error and terminate.
func (a *App) initDatabase() {
	maxOpen := a.config.GetInt(`database.max.open`)
	maxIdle := a.config.GetInt(`database.max.idle`)
	maxLifetime := a.config.GetInt(`database.max.lifetime`)
	maxIdleTime := a.config.GetInt(`database.max.idletime`)

	dsn := a.dsnMySQL()
	driver := dbops.MySQLDriver
	queryBuilder := goqu.Dialect(dbops.MySQLDriver)
	if a.config.GetString(`database.driver`) == dbops.PostgresDriver {
		dsn = a.dsnPostgreSQL()
		driver = dbops.PostgresDriver
		queryBuilder = goqu.Dialect(dbops.PostgresDriver)
	}

	database, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalln("failed to open database", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatalln("failed to ping database connection", err)
	}

	database.SetMaxOpenConns(int(maxOpen))
	database.SetMaxIdleConns(int(maxIdle))
	database.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
	database.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)

	dbops.SetVerbose(true)
	a.database = database
	a.queryBuilder = queryBuilder
	a.transaction = dbops.NewTransaction(database)
}

// initRedis initializes a Redis client using settings from the configuration.
// It verifies the connection by pinging the Redis server. If the connection fails,
// the application will log a fatal error and terminate.
func (a *App) initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr: a.config.GetString("redis.addr"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("failed to init redis", err)
	}

	a.redisDB = rdb
}

func (a *App) initMessaging() {
	// msg, err := redispubsub.NewClient(
	// 	"",
	// 	redispubsub.WithExistingClient(a.redisDB),
	// 	redispubsub.WithLogger(a.telemetry.Logger()),
	// 	redispubsub.WithSyncPublisher(),
	// )

	msg, err := googlepubsub.NewClient(
		context.Background(),
		a.config.GetString("pubsub.project.id"),
		googlepubsub.WithLogger(a.telemetry.Logger()),
		googlepubsub.WithAutoAck(),
	)
	if err != nil {
		log.Fatalln("failed to init messaging", err)
	}

	a.messaging = msg
}

// initHTTPServer initializes the HTTP server with settings from the configuration,
// such as address, timeouts, and the router handler wrapped with CORS middleware.
// This method should be called after initializing the router to ensure the server
// is ready to handle incoming requests.
func (a *App) initHTTPServer() {
	a.httpRouter = framework.NewRouter()
	a.httpServer = &http.Server{
		Addr: a.config.GetString("server.address.http"),
		Handler: framework.Chain(
			a.httpRouter,
			framework.Recovery,
			// next-mr: Need to create a custom CORS implementation to standardize error messages
			cors.AllowAll().Handler,
			instrument.UseTelemetryServer(a.telemetry, a.uuid.Generate),
			framework.JWT(a.jwt, "gostarter.access.token", "/auth", "/graphql/playground"),
		),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (a *App) initGQLServer() {
	a.gqlRouter = framework.NewRouter()
	a.gqlServer = &http.Server{
		Addr: a.config.GetString("server.address.gql"),
		Handler: framework.Chain(
			a.gqlRouter,
			framework.Recovery,
			cors.AllowAll().Handler,
			instrument.UseTelemetryServer(a.telemetry, a.uuid.Generate),
			framework.JWT(a.jwt, "gostarter.access.token", "/graphql/playground"),
		),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (a *App) initGRPCServer() {
	opts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(framework.UnaryServerRecovery)}
	opts = append(opts, instrument.UnaryTelemetryServerInterceptor(a.telemetry, a.uuid.Generate)...)
	opts = append(opts, grpc.ChainUnaryInterceptor(
		framework.UnaryServerError,
		framework.UnaryServerJWT(a.jwt, "gostarter.access.token", "/gostarter.api.auth.AuthService"),
		framework.UnaryServerProtoValidate(a.protoValidator),
	))

	server := grpc.NewServer(opts...)
	reflection.Register(server)
	a.grpcServer = server
}

// initTasks starts all background tasks or services registered with the application.
// If any task fails to start, the application will log a fatal error and terminate.
func (a *App) initTasks() {
	for _, runnable := range a.runnables {
		if err := runnable.Start(); err != nil {
			log.Fatalln("failed to init initTasks", err)
		}
	}
}

// initClosers registers cleanup functions for all core components of the application.
// These functions are responsible for gracefully shutting down the HTTP server,
// gRPC Server, closing database connections, terminating the Redis client, and
// cleaning up the configuration. This method is typically called when stopping
// the application to ensure all resources are released properly.
func (a *App) initClosers() {
	a.closerFn = map[string]func(context.Context) error{
		"HTTP Server": func(ctx context.Context) error {
			return a.httpServer.Shutdown(ctx)
		},
		"GQL Server": func(ctx context.Context) error {
			return a.gqlServer.Shutdown(ctx)
		},
		"GRPC Server": func(_ context.Context) error { //nolint:unparam // its ok
			a.grpcServer.GracefulStop()

			return nil
		},
		"Database": func(_ context.Context) error {
			return a.database.Close()
		},
		"Redis": func(_ context.Context) error {
			return a.redisDB.Close()
		},
		"Config": func(_ context.Context) error {
			return a.config.Close()
		},
		"Telemetry": func(_ context.Context) error {
			return a.telemetry.Close()
		},
	}
}
