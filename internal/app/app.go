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
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/codec"
	"github.com/shandysiswandi/goreng/config"
	"github.com/shandysiswandi/goreng/goroutine"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/jwt"
	"github.com/shandysiswandi/goreng/messaging"
	"github.com/shandysiswandi/goreng/task"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
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
	sqlkitDB       *sqlkit.DB
	redisDB        *redis.Client
	messaging      messaging.Client
	httpServer     *http.Server
	gqlServer      *http.Server
	grpcServer     *grpc.Server
	httpRouter     *framework.Router
	gqlRouter      *framework.Router
	goroutine      *goroutine.Manager
	hash           hash.Hash
	secHash        hash.Hash
	jwt            jwt.JWT
	clock          clock.Clocker
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
	app.initJWT()
	app.initLibraries()
	app.initDatabase()
	app.initRedis()
	app.initMessaging()
	app.initHTTPServer()
	app.initGQLServer()
	app.initGRPCServer()
	app.initModules()
	app.initTasks()
	app.initClosers()

	return app
}
