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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// initSTDLog initializes the standard logger with specific flags to include the date, time,
// and file location in log messages. This method is typically called during the initialization
// phase of the application to ensure consistent logging behavior.
func (a *App) initSTDLog() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// initConfig initializes the application's configuration by loading settings from a specified
// YAML file using the Viper library. If the configuration cannot be loaded, the application
// will log a fatal error and terminate. This method should be called early in the application's
// initialization process.
func (a *App) initConfig() {
	cfg, err := config.NewViperConfig("config/config.yaml")
	if err != nil {
		log.Fatalln("failed to init config", err)
	}

	a.config = cfg
}

// initDatabase initializes the application's database connection using settings from the
// configuration. It sets up connection pooling parameters and tests the connection by pinging
// the database. If any step fails, the application will log a fatal error and terminate.
func (a *App) initDatabase() {
	maxOpen := a.config.GetInt(`database.max.open`)
	maxIdle := a.config.GetInt(`database.max.idle`)
	maxLifetime := a.config.GetInt(`database.max.lifetime`)
	maxIdletime := a.config.GetInt(`database.max.idletime`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		a.config.GetString(`database.user`),
		a.config.GetString(`database.pass`),
		a.config.GetString(`database.host`),
		a.config.GetString(`database.port`),
		a.config.GetString(`database.name`),
	)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	val.Encode()

	database, err := sql.Open("mysql", fmt.Sprintf("%s?%s", dsn, val.Encode()))
	if err != nil {
		log.Fatalln("failed to open database", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatalln("failed to ping database connection", err)
	}

	database.SetMaxOpenConns(int(maxOpen))
	database.SetMaxIdleConns(int(maxIdle))
	database.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
	database.SetConnMaxIdleTime(time.Duration(maxIdletime) * time.Minute)

	a.database = database
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

	a.redisdb = rdb
}

// initHTTPRouter initializes the HTTP router using the httprouter library.
// It sets up custom handlers for "Not Found" and "Method Not Allowed" errors,
// returning JSON responses with appropriate status codes.
func (a *App) initHTTPRouter() {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{"code": "40400", "message": "Endpoint not found"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(map[string]string{"code": "40500", "message": "Method endpoint not allowed"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	a.httpRouter = router
}

// initHTTPServer initializes the HTTP server with settings from the configuration,
// such as address, timeouts, and the router handler wrapped with CORS middleware.
// This method should be called after initializing the router to ensure the server
// is ready to handle incoming requests.
func (a *App) initHTTPServer() {
	handler := cors.Default().Handler(a.httpRouter)

	a.httpServer = &http.Server{
		Addr:              a.config.GetString("server.address.http"),
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (a *App) initGRPCServer() {
	server := grpc.NewServer()
	reflection.Register(server)
	a.grpcServer = server
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
		log.Fatalln("failed to init validation protovalidate", err)
	}

	a.uidnumber = snow
	a.pvalidator = pvalidator
	a.clock = clock.NewClock()
	a.uuid = uid.NewUUIDString()
	a.codecJSON = codec.NewJSONCodec()
	a.codecMsgPack = codec.NewMsgpackCodec()
	a.validator = validation.NewV10Validator()
	a.logger = logger.NewStdLogger()
}

// initTasks starts all background tasks or services registered with the application.
// If any task fails to start, the application will log a fatal error and terminate.
func (a *App) initTasks() {
	for _, runnable := range a.runables {
		if err := runnable.Start(); err != nil {
			log.Fatalln("failed to init initTasks", err)
		}
	}
}

// iniClosed registers cleanup functions for all core components of the application.
// These functions are responsible for gracefully shutting down the HTTP server,
// gRPC Server, closing database connections, terminating the Redis client, and
// cleaning up the configuration. This method is typically called when stopping
// the application to ensure all resources are released properly.
func (a *App) iniClosed() {
	a.closersFn = append(a.closersFn, []func(ctx context.Context) error{
		func(ctx context.Context) error {
			return a.httpServer.Shutdown(ctx)
		},
		func(_ context.Context) error {
			a.grpcServer.GracefulStop()

			return nil
		},
		func(context.Context) error {
			return a.database.Close()
		},
		func(context.Context) error {
			return a.redisdb.Close()
		},
		func(context.Context) error {
			return a.config.Close()
		},
	}...)
}
