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

	"github.com/doug-martin/goqu/v9"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/goroutine"
	"github.com/shandysiswandi/gostarter/pkg/hash"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/task"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
)

// App encapsulates the core components of the application, including configuration,
// UID generators, clock, codecs, database connections, Redis client, HTTP server,
// HTTP router, and tasks. It manages the initialization and graceful shutdown of
// these components.
type App struct {
	config         config.Config
	uidNumber      uid.NumberID
	uuid           uid.StringID
	codecJSON      codec.Codec
	codecMsgPack   codec.Codec
	validator      validation.Validator
	protoValidator validation.Validator
	telemetry      *telemetry.Telemetry
	database       *sql.DB
	transaction    dbops.Tx
	queryBuilder   goqu.DialectWrapper
	redisDB        *redis.Client
	messaging      messaging.Client
	httpServer     *http.Server
	grpcServer     *grpc.Server
	httpRouter     *httpserver.Router
	goroutine      *goroutine.Manager
	hash           hash.Hash
	secHash        hash.Hash
	jwt            jwt.JWT
	runnables      []task.Runner
	closerFn       map[string]func(context.Context) error
}

// New creates and returns a new instance of the App structure. This function initializes
// all the core components of the application, including standard logging, configuration,
// database connections, Redis client, HTTP router, HTTP server, gRPC Server and various
// utility libraries. This method is typically called before starting the application to
// ensure that all components are properly set up.
func New() *App {
	app := &App{}

	app.initConfig()
	app.initTelemetry()
	app.initLibraries()
	app.initDatabase()
	app.initRedis()
	app.initMessaging()
	app.initHTTPServer()
	app.initGRPCServer()
	app.initModules()
	app.initTasks()
	app.initClosers()

	return app
}
