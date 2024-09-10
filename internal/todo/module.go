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
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/task"
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
	Router         *httprouter.Router
	GRPCServer     *grpc.Server
	Logger         logger.Logger
}

func New(dep Dependency) (*Expose, error) {
	// init outbound | database | http client | grpc client | redis | etc.
	sqlTodo := outbound.NewSQLTodo(dep.Database, dep.Config)

	// init services | useCases | business logic
	findUC := service.NewFind(dep.Logger, sqlTodo, dep.Validator)
	fetchUC := service.NewFetch(dep.Logger, sqlTodo)
	createUC := service.NewCreate(dep.Logger, sqlTodo, dep.Validator, dep.UIDNumber)
	deleteUC := service.NewDelete(dep.Logger, sqlTodo, dep.Validator)
	updateUC := service.NewUpdate(dep.Logger, sqlTodo, dep.Validator)
	updateStatusUC := service.NewUpdateStatus(dep.Logger, sqlTodo, dep.Validator)

	// register endpoint REST
	inboundhttp.RegisterRESTEndpoint(dep.Router, &inboundhttp.Endpoint{
		FindUC:         findUC,
		CreateUC:       createUC,
		DeleteUC:       deleteUC,
		FetchUC:        fetchUC,
		UpdateUC:       updateUC,
		UpdateStatusUC: updateStatusUC,
	})

	inboundhttp.RegisterSSEEndpoint(dep.Router, inboundhttp.NewSSE(
		dep.CodecJSON,
		dep.Logger,
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
