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
	"github.com/shandysiswandi/gostarter/internal/region"
	"github.com/shandysiswandi/gostarter/internal/shortly"
	"github.com/shandysiswandi/gostarter/internal/todo"
)

func (a *App) initModules() {
	a.moduleAuth()
	a.moduleShortly()
	a.moduleRegion()
	a.moduleTodo()
}

func (a *App) moduleAuth() {
	if a.config.GetBool("module.flag.auth") {
		_, err := auth.New(auth.Dependency{
			Config:    a.config,
			Database:  a.database,
			Telemetry: a.telemetry,
			Router:    a.httpRouter,
			Validator: a.validator,
			UIDNumber: a.uidNumber,
			Hash:      a.hash,
			SecHash:   a.secHash,
			JWT:       a.jwt,
		})
		if err != nil {
			log.Fatalln("failed to init module auth", err)
		}
	}
}

func (a *App) moduleShortly() {
	if a.config.GetBool("module.flag.shortly") {
		_, err := shortly.New(shortly.Dependency{
			RedisDB:   a.redisDB,
			Config:    a.config,
			CodecJSON: a.codecJSON,
			Validator: a.validator,
			Router:    a.httpRouter,
			Telemetry: a.telemetry,
		})
		if err != nil {
			log.Fatalln("failed to init module shortly", err)
		}
	}
}

func (a *App) moduleRegion() {
	if a.config.GetBool("module.flag.region") {
		_, err := region.New(region.Dependency{
			Database:  a.database,
			RedisDB:   a.redisDB,
			Config:    a.config,
			CodecJSON: a.codecJSON,
			Validator: a.validator,
			Router:    a.httpRouter,
			Telemetry: a.telemetry,
		})
		if err != nil {
			log.Fatalln("failed to init module region", err)
		}
	}
}

func (a *App) moduleTodo() {
	if a.config.GetBool("module.flag.todo") {
		expTodo, err := todo.New(todo.Dependency{
			Database:       a.database,
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
