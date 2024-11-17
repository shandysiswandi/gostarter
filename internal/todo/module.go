package todo

import (
	"database/sql"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/doug-martin/goqu/v9"
	"github.com/redis/go-redis/v9"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	inboundgql "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/gql"
	inboundgrpc "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/grpc"
	inboundhttp "github.com/shandysiswandi/gostarter/internal/todo/internal/inbound/http"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/job"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/framework/gql"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/framework/middleware"
	"github.com/shandysiswandi/gostarter/pkg/goroutine"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
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
	Database       *sql.DB
	QueryBuilder   goqu.DialectWrapper
	RedisDB        *redis.Client
	Messaging      messaging.Client
	Config         config.Config
	UIDNumber      uid.NumberID
	CodecJSON      codec.Codec
	Validator      validation.Validator
	ProtoValidator validation.Validator
	JWT            jwt.JWT
	Router         *httpserver.Router
	GRPCServer     *grpc.Server
	Telemetry      *telemetry.Telemetry
	Goroutine      *goroutine.Manager
}

const todoPath = "/todos"

//nolint:funlen // it's long line because it format param dependency
func New(dep Dependency) (*Expose, error) {
	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes outbound dependencies for core services.
	// This includes setups for outbound services: Database, HTTP client, gRPC client, Redis, etc.
	sqlTodo := outbound.NewSQLTodo(dep.Database, dep.QueryBuilder)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes core business logic or use cases to handle user interaction
	findUC := service.NewFind(dep.Telemetry, sqlTodo, dep.Validator)
	fetchUC := service.NewFetch(dep.Telemetry, sqlTodo)
	createUC := service.NewCreate(dep.Telemetry, sqlTodo, dep.Validator, dep.UIDNumber)
	deleteUC := service.NewDelete(dep.Telemetry, sqlTodo, dep.Validator)
	updateUC := service.NewUpdate(dep.Telemetry, sqlTodo, dep.Validator)
	updateStatusUC := service.NewUpdateStatus(dep.Telemetry, sqlTodo, dep.Validator)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes REST API endpoints to handle core user workflows:
	eHTTP := inboundhttp.NewEndpoint(createUC, deleteUC, findUC, fetchUC,
		updateStatusUC, updateUC)
	eSSE := inboundhttp.NewSSE(dep.CodecJSON)

	//
	dep.Router.Endpoint(http.MethodGet, todoPath+"/:id", eHTTP.Find)
	dep.Router.Endpoint(http.MethodGet, todoPath, eHTTP.Fetch)
	dep.Router.Endpoint(http.MethodPost, todoPath, eHTTP.Create)
	dep.Router.Endpoint(http.MethodPut, todoPath+"/:id", eHTTP.Update)
	dep.Router.Endpoint(http.MethodPatch, todoPath+"/:id/status", eHTTP.UpdateStatus)
	dep.Router.Endpoint(http.MethodDelete, todoPath+"/:id", eHTTP.Delete)
	//
	dep.Router.Native(http.MethodGet, "/events", http.HandlerFunc(eSSE.HandleEvent), middleware.Recovery)
	dep.Router.Native(http.MethodGet, "/trigger-event", http.HandlerFunc(eSSE.HandleEvent),
		middleware.Recovery)
	//

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes gRPC API endpoints to handle core user workflows:
	eGRPC := inboundgrpc.NewEndpoint(dep.ProtoValidator, findUC, fetchUC,
		createUC, deleteUC, updateUC, updateStatusUC)
	pb.RegisterTodoServiceServer(dep.GRPCServer, eGRPC)

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes graphQL API endpoints to handle core user workflows:
	eGQL := &inboundgql.Endpoint{
		FindUC:         findUC,
		CreateUC:       createUC,
		DeleteUC:       deleteUC,
		FetchUC:        fetchUC,
		UpdateUC:       updateUC,
		UpdateStatusUC: updateStatusUC,
	}
	exec := ql.NewExecutableSchema(ql.Config{Resolvers: eGQL})
	gqlServer := gql.ServerDefault(exec)
	dep.Router.Native(http.MethodPost, "/graphql", gqlServer)
	if dep.Config.GetBool("feature.flag.graphql.playground") {
		dep.Router.Native(http.MethodGet, "/graphql/playground",
			playground.Handler("GraphQL playground", "/graphql"))
	}

	// ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== ===== =====
	// This block initializes runner job to handle background workflows:
	todoSub := &job.TodoSubscriber{
		MsgClient: dep.Messaging,
		Tel:       dep.Telemetry,
	}
	todoPub := &job.TodoPublisher{
		MsgClient: dep.Messaging,
		Tel:       dep.Telemetry,
	}

	return &Expose{
		Tasks: []task.Runner{
			todoSub,
			todoPub,
		},
	}, nil
}
