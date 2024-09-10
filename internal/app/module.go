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

	"github.com/shandysiswandi/gostarter/internal/gallery"
	"github.com/shandysiswandi/gostarter/internal/region"
	"github.com/shandysiswandi/gostarter/internal/shortly"
	"github.com/shandysiswandi/gostarter/internal/todo"
)

func (a *App) initModules() {
	if a.config.GetBool("module.flag.gallery") {
		_, err := gallery.New(gallery.Dependency{
			RedisDB:   a.redisDB,
			Config:    a.config,
			CodecJSON: a.codecJSON,
			Validator: a.validator,
			Router:    a.httpRouter,
			Logger:    a.logger,
		})
		if err != nil {
			log.Fatalln("failed to init module gallery", err)
		}
	}

	if a.config.GetBool("module.flag.shortly") {
		_, err := shortly.New(shortly.Dependency{
			RedisDB:   a.redisDB,
			Config:    a.config,
			CodecJSON: a.codecJSON,
			Validator: a.validator,
			Router:    a.httpRouter,
			Logger:    a.logger,
		})
		if err != nil {
			log.Fatalln("failed to init module shortly", err)
		}
	}

	if a.config.GetBool("module.flag.region") {
		_, err := region.New(region.Dependency{
			Database:  a.database,
			RedisDB:   a.redisDB,
			Config:    a.config,
			CodecJSON: a.codecJSON,
			Validator: a.validator,
			Router:    a.httpRouter,
			Logger:    a.logger,
		})
		if err != nil {
			log.Fatalln("failed to init module region", err)
		}
	}

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
			Logger:         a.logger,
		})
		if err != nil {
			log.Fatalln("failed to init module todo", err)
		}

		a.runnables = append(a.runnables, expTodo.Tasks...)
	}
}
