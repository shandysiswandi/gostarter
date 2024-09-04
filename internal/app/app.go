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
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/task"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
)

// App encapsulates the core components of the application, including configuration,
// UID generators, clock, codecs, database connections, Redis client, HTTP server,
// HTTP router, and tasks. It manages the initialization and graceful shutdown of
// these components.
type App struct {
	config       config.Config
	uidnumber    uid.NumberID
	uuid         uid.StringID
	codecJSON    codec.Codec
	codecMsgPack codec.Codec
	validator    validation.Validator
	pvalidator   validation.Validator
	logger       logger.Logger
	database     *sql.DB
	redisdb      *redis.Client
	httpServer   *http.Server
	grpcServer   *grpc.Server
	httpRouter   *httprouter.Router
	runables     []task.Runner
	closersFn    []func(context.Context) error
}

// New creates and returns a new instance of the App structure. This function initializes
// all the core components of the application, including standard logging, configuration,
// database connections, Redis client, HTTP router, HTTP server, gRPC Srever and various
// utility libraries. This method is typically called before starting the application to
// ensure that all components are properly set up.
func New() *App {
	app := &App{}

	app.initSTDLog()
	app.initConfig()
	app.initDatabase()
	app.initRedis()
	app.initHTTPRouter()
	app.initHTTPServer()
	app.initGRPCServer()
	app.initLibraries()
	app.initModules()
	app.initTasks()
	app.iniClosed()

	return app
}
