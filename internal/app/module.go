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

	"github.com/shandysiswandi/gostarter/internal/region"
	"github.com/shandysiswandi/gostarter/internal/shortly"
	"github.com/shandysiswandi/gostarter/internal/todo"
)

func (a *App) initModules() {
	expTodo, err := todo.New(todo.Dependency{
		Database:       a.database,
		RedisDB:        a.redisdb,
		Config:         a.config,
		UIDNumber:      a.uidnumber,
		CodecJSON:      a.codecJSON,
		Validator:      a.validator,
		ProtoValidator: a.pvalidator,
		Router:         a.httpRouter,
		GRPCServer:     a.grpcServer,
		Logger:         a.logger,
	})
	if err != nil {
		log.Fatalln("failed to init module todo", err)
	}

	_, err = region.New(region.Dependency{
		Database:  a.database,
		RedisDB:   a.redisdb,
		Config:    a.config,
		CodecJSON: a.codecJSON,
		Validator: a.validator,
		Router:    a.httpRouter,
		Logger:    a.logger,
	})
	if err != nil {
		log.Fatalln("failed to init module region", err)
	}

	_, err = shortly.New(shortly.Dependency{
		RedisDB:   a.redisdb,
		Config:    a.config,
		CodecJSON: a.codecJSON,
		Validator: a.validator,
		Router:    a.httpRouter,
		Logger:    a.logger,
	})
	if err != nil {
		log.Fatalln("failed to init module shortly", err)
	}

	a.runables = append(a.runables, expTodo.Tasks...)
}
