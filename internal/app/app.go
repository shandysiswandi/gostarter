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
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/task"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

// App encapsulates the core components of the application, including configuration,
// UID generators, clock, codecs, database connections, Redis client, HTTP server,
// HTTP router, and tasks. It manages the initialization and graceful shutdown of
// these components.
type App struct {
	config       config.Config
	uidnumber    uid.Number
	uuid         uid.String
	clock        clock.Clocker
	codecJSON    codec.Codec
	codecMsgPack codec.Codec
	validator    validation.Validator
	database     *sql.DB
	redisdb      *redis.Client
	httpServer   *http.Server
	httpRouter   *httprouter.Router
	runables     []task.Runner
	closersFn    []func(context.Context) error
}

// New creates and returns a new instance of the App structure. This function initializes
// an empty App object, which can then be configured and started by calling its methods.
func New() *App {
	return &App{}
}

// ensureInitialized initializes all the core components of the application, including
// standard logging, configuration, database connections, Redis client, HTTP router,
// HTTP server, and various utility libraries. This method is typically called before
// starting the application to ensure that all components are properly set up.
func (a *App) ensureInitialized() {
	a.initSTDLog()
	a.initConfig()
	a.initDatabase()
	a.initRedis()
	a.initHTTPRouter()
	a.initHTTPServer()
	a.initLibraries()
	a.initModules()
	a.initTasks()
}

// ensureClosed registers cleanup functions for all core components of the application.
// These functions are responsible for gracefully shutting down the HTTP server,
// closing database connections, terminating the Redis client, and cleaning up
// the configuration. This method is typically called when stopping the application
// to ensure all resources are released properly.
func (a *App) ensureClosed() {
	a.closersFn = append(a.closersFn, []func(ctx context.Context) error{
		func(ctx context.Context) error {
			return a.httpServer.Shutdown(ctx)
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
