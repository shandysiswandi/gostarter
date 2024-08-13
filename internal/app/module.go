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

	"github.com/shandysiswandi/gostarter/internal/todo"
)

func (a *App) initModules() {
	expTodo, err := todo.New(todo.Dependency{
		Database:       a.database,
		RedisDB:        a.redisdb,
		Config:         a.config,
		UIDNumber:      a.uidnumber,
		Clock:          a.clock,
		CodecJSON:      a.codecJSON,
		Validator:      a.validator,
		ProtoValidator: a.pvalidator,
		Router:         a.httpRouter,
		GRPCServer:     a.grpcServer,
	})
	if err != nil {
		log.Fatalln("failed to init module todo", err)
	}

	a.runables = append(a.runables, expTodo.Tasks...)
}
