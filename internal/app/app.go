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

func New() *App {
	return &App{}
}

func (a *App) ensureInitialized() {
	a.initSTDLog()
	a.initConfig()
	a.initDatabase()
	a.initRedis()
	a.initHTTPRouter()
	a.initHTTPServer()
	a.initLibraries()
	a.initTasks()
	a.initModules()
}

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
