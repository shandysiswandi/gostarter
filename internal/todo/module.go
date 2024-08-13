package todo

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	pb "github.com/shandysiswandi/gostarter/api/gen-go/todo"
	inboundgrpc "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/grpc"
	inboundhttp "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/http"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/job"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
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
	UIDNumber      uid.Number
	Clock          clock.Clocker
	CodecJSON      codec.Codec
	Validator      validation.Validator
	ProtoValidator validation.Validator
	Router         *httprouter.Router
	GRPCServer     *grpc.Server
}

func New(dep Dependency) (*Expose, error) {
	// init outbound | database | http client | grpc client | redis | etc.
	mysqlTodo := outbound.NewMysqlTodo(dep.Database)

	// init service | usecase | interactor | logic
	getByIDUC := service.NewGetByID(mysqlTodo, dep.Validator)
	getWithFilterUC := service.NewGetWithFilter(mysqlTodo, dep.Validator)
	createUC := service.NewCreate(mysqlTodo, dep.Validator, dep.UIDNumber)
	deleteUC := service.NewDelete(mysqlTodo, dep.Validator)
	updateUC := service.NewUpdate(mysqlTodo, dep.Validator)
	updateStatusUC := service.NewUpdateStatus(mysqlTodo, dep.Validator)

	// init inbound endpoint http
	endpointHTTP := &inboundhttp.Endpoint{
		GetByIDUC:       getByIDUC,
		CreateUC:        createUC,
		DeleteUC:        deleteUC,
		GetWithFilterUC: getWithFilterUC,
		UpdateUC:        updateUC,
		UpdateStatusUC:  updateStatusUC,
	}

	// init inbound endpoint grpc
	endpointGRPC := inboundgrpc.NewEndpoint(
		dep.ProtoValidator,
		getByIDUC,
		getWithFilterUC,
		createUC,
		deleteUC,
		updateUC,
		updateStatusUC,
	)

	// register endpoint HTTP
	inboundhttp.RegisterHTTPEndpoint(dep.Router, endpointHTTP)

	// register endpoint GRPC | GRAPHQL | etc.
	pb.RegisterTodoServiceServer(dep.GRPCServer, endpointGRPC)

	// jobs | background tasks
	exampleJob := &job.ExampleJob{}

	return &Expose{
		Tasks: []task.Runner{exampleJob},
	}, nil
}
