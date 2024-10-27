package todo

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	inboundgql "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/gql"
	inboundgrpc "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/grpc"
	inboundhttp "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/http"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/job"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/goroutine"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
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
	Database       *sql.DB
	RedisDB        *redis.Client
	Config         config.Config
	UIDNumber      uid.NumberID
	CodecJSON      codec.Codec
	Validator      validation.Validator
	ProtoValidator validation.Validator
	JWT            jwt.JWT
	Router         *httprouter.Router
	GRPCServer     *grpc.Server
	Telemetry      *telemetry.Telemetry
	Goroutine      *goroutine.Manager
}

func New(dep Dependency) (*Expose, error) {
	// init outbound | database | http client | grpc client | redis | etc.
	sqlTodo := outbound.NewSQLTodo(dep.Database, dep.Config)

	// init services | useCases | business logic
	findUC := service.NewFind(dep.Telemetry, sqlTodo, dep.Validator)
	fetchUC := service.NewFetch(dep.Telemetry, sqlTodo)
	createUC := service.NewCreate(dep.Telemetry, sqlTodo, dep.Validator, dep.UIDNumber)
	deleteUC := service.NewDelete(dep.Telemetry, sqlTodo, dep.Validator)
	updateUC := service.NewUpdate(dep.Telemetry, sqlTodo, dep.Validator)
	updateStatusUC := service.NewUpdateStatus(dep.Telemetry, sqlTodo, dep.Validator)

	// register endpoint REST
	endpoint := inboundhttp.NewEndpoint(
		createUC,
		deleteUC,
		findUC,
		fetchUC,
		updateStatusUC,
		updateUC,
		dep.Goroutine,
	)
	inboundhttp.RegisterRESTEndpoint(dep.Router, endpoint, dep.JWT)
	inboundhttp.RegisterSSEEndpoint(dep.Router, inboundhttp.NewSSE(
		dep.CodecJSON,
	))

	// register endpoint GRPC
	pb.RegisterTodoServiceServer(dep.GRPCServer, inboundgrpc.NewEndpoint(
		dep.ProtoValidator,
		findUC,
		fetchUC,
		createUC,
		deleteUC,
		updateUC,
		updateStatusUC,
	))

	// register endpoint GRAPHQL
	inboundgql.RegisterGQLEndpoint(dep.Router, dep.Config, &inboundgql.Endpoint{
		FindUC:         findUC,
		CreateUC:       createUC,
		DeleteUC:       deleteUC,
		FetchUC:        fetchUC,
		UpdateUC:       updateUC,
		UpdateStatusUC: updateStatusUC,
	})

	// jobs | background tasks
	exampleJob := &job.ExampleJob{}

	return &Expose{
		Tasks: []task.Runner{exampleJob},
	}, nil
}
