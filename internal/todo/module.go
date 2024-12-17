package todo

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/job"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/task"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	"google.golang.org/grpc"
)

type Expose struct {
	Tasks []task.Runner
}

type Dependency struct {
	Database     *sql.DB
	QueryBuilder goqu.DialectWrapper
	Messaging    messaging.Client
	Config       config.Config
	UIDNumber    uid.NumberID
	CodecJSON    codec.Codec
	Validator    validation.Validator
	Router       *framework.Router
	GQLRouter    *framework.Router
	GRPCServer   *grpc.Server
	Telemetry    *telemetry.Telemetry
}

func New(dep Dependency) (*Expose, error) {
	// This block initializes outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlTodo := outbound.NewSQLTodo(dep.Database, dep.QueryBuilder, dep.Telemetry)

	// This block initializes core business logic or use cases to handle user interaction
	ucDep := usecase.Dependency{
		UIDNumber: dep.UIDNumber,
		Validator: dep.Validator,
		Telemetry: dep.Telemetry,
	}
	findUC := usecase.NewFind(ucDep, sqlTodo)
	fetchUC := usecase.NewFetch(ucDep, sqlTodo)
	createUC := usecase.NewCreate(ucDep, sqlTodo)
	deleteUC := usecase.NewDelete(ucDep, sqlTodo)
	updateUC := usecase.NewUpdate(ucDep, sqlTodo)
	updateStatusUC := usecase.NewUpdateStatus(ucDep, sqlTodo)

	// This block initializes REST, SSE, gRPC, and graphQL API endpoints to handle core user workflows:
	inbound := inbound.Inbound{
		Router:     dep.Router,
		GQLRouter:  dep.GQLRouter,
		GRPCServer: dep.GRPCServer,
		CodecJSON:  dep.CodecJSON,
		Telemetry:  dep.Telemetry,
		//
		CreateUC:       createUC,
		DeleteUC:       deleteUC,
		FindUC:         findUC,
		FetchUC:        fetchUC,
		UpdateStatusUC: updateStatusUC,
		UpdateUC:       updateUC,
	}
	inbound.RegisterTodoServiceServer()

	// This block initializes runner job to handle background workflows:
	jobs := job.New(job.Dependency{
		Messaging:    dep.Messaging,
		Config:       dep.Config,
		CodecJSON:    dep.CodecJSON,
		Telemetry:    dep.Telemetry,
		DomainCreate: createUC,
	})

	return &Expose{Tasks: jobs}, nil
}
