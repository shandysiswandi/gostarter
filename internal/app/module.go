// Package app provides a structured framework for initializing, configuring,
// and managing the core components of the application. It handles the setup
// of essential services such as logging, configuration management, database
// connections, Redis client, HTTP router and server, utility libraries, and
// background tasks. The package ensures that all components are properly
// initialized at startup and gracefully shut down when the application stops,
// offering a robust foundation for building and maintaining the application.
package app

import (
	"log"

	"github.com/shandysiswandi/gostarter/internal/auth"
	"github.com/shandysiswandi/gostarter/internal/todo"
)

func (a *App) initModules() {
	a.moduleAuth()
	a.moduleTodo()
}

func (a *App) moduleAuth() {
	if a.config.GetBool("module.flag.auth") {
		_, err := auth.New(auth.Dependency{
			Database:       a.database,
			QueryBuilder:   a.queryBuilder,
			Telemetry:      a.telemetry,
			Router:         a.httpRouter,
			GRPCServer:     a.grpcServer,
			ProtoValidator: a.protoValidator,
			Validator:      a.validator,
			UIDNumber:      a.uidNumber,
			Hash:           a.hash,
			SecHash:        a.secHash,
			JWT:            a.jwt,
		})
		if err != nil {
			log.Fatalln("failed to init module auth", err)
		}
	}
}

func (a *App) moduleTodo() {
	if a.config.GetBool("module.flag.todo") {
		expTodo, err := todo.New(todo.Dependency{
			Database:       a.database,
			QueryBuilder:   a.queryBuilder,
			RedisDB:        a.redisDB,
			Config:         a.config,
			UIDNumber:      a.uidNumber,
			CodecJSON:      a.codecJSON,
			Validator:      a.validator,
			ProtoValidator: a.protoValidator,
			Router:         a.httpRouter,
			GRPCServer:     a.grpcServer,
			Telemetry:      a.telemetry,
			Goroutine:      a.goroutine,
			JWT:            a.jwt,
		})
		if err != nil {
			log.Fatalln("failed to init module todo", err)
		}

		a.runnables = append(a.runnables, expTodo.Tasks...)
	}
}
