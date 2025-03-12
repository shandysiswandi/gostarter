package app

import (
	"log"

	"github.com/shandysiswandi/gostarter/internal/auth"
	"github.com/shandysiswandi/gostarter/internal/payment"
	"github.com/shandysiswandi/gostarter/internal/rbac"
	"github.com/shandysiswandi/gostarter/internal/todo"
	"github.com/shandysiswandi/gostarter/internal/user"
)

func (a *App) initModules() {
	a.moduleAuth()
	a.moduleRBAC()
	a.modulePayment()
	a.moduleTodo()
	a.moduleUser()
}

func (a *App) moduleAuth() {
	if a.config.GetBool("module.flag.auth") {
		_, err := auth.New(auth.Dependency{
			Telemetry:  a.telemetry,
			SQLKitDB:   a.sqlkitDB,
			Router:     a.httpRouter,
			GRPCServer: a.grpcServer,
			Validator:  a.validator,
			UIDNumber:  a.uidNumber,
			Hash:       a.hash,
			SecHash:    a.secHash,
			JWT:        a.jwt,
			Clock:      a.clock,
		})
		if err != nil {
			log.Fatalln("failed to init module auth", err)
		}
	}
}

func (a *App) moduleRBAC() {
	if a.config.GetBool("module.flag.rbac") {
		_, err := rbac.New(rbac.Dependency{
			Telemetry: a.telemetry,
			Router:    a.httpRouter,
			Validator: a.validator,
			UIDNumber: a.uidNumber,
			Clock:     a.clock,
		})
		if err != nil {
			log.Fatalln("failed to init module rbac", err)
		}
	}
}

func (a *App) modulePayment() {
	if a.config.GetBool("module.flag.payment") {
		_, err := payment.New(payment.Dependency{
			UIDNumber: a.uidNumber,
			Validator: a.validator,
			Router:    a.httpRouter,
			Telemetry: a.telemetry,
			Hash:      nil,
			SecHash:   nil,
			Clock:     a.clock,
		})
		if err != nil {
			log.Fatalln("failed to init module payment", err)
		}
	}
}

func (a *App) moduleTodo() {
	if a.config.GetBool("module.flag.todo") {
		expTodo, err := todo.New(todo.Dependency{
			Messaging:  a.messaging,
			Config:     a.config,
			UIDNumber:  a.uidNumber,
			CodecJSON:  a.codecJSON,
			Validator:  a.validator,
			Router:     a.httpRouter,
			GQLRouter:  a.gqlRouter,
			GRPCServer: a.grpcServer,
			Telemetry:  a.telemetry,
		})
		if err != nil {
			log.Fatalln("failed to init module todo", err)
		}

		a.runnables = append(a.runnables, expTodo.Tasks...)
	}
}

func (a *App) moduleUser() {
	if a.config.GetBool("module.flag.user") {
		_, err := user.New(user.Dependency{
			Validator: a.validator,
			Hash:      a.hash,
			Router:    a.httpRouter,
			Telemetry: a.telemetry,
		})
		if err != nil {
			log.Fatalln("failed to init module user", err)
		}
	}
}
