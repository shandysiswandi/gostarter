package inbound

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"google.golang.org/grpc"
)

var (
	errFailedParseToUint = goerror.NewInvalidFormat("failed parse id to uint")
	errInvalidBody       = goerror.NewInvalidFormat("invalid request body")
)

type Inbound struct {
	Config     config.Config
	Router     *framework.Router
	GQLRouter  *framework.Router
	GRPCServer *grpc.Server
	CodecJSON  codec.Codec
	//
	CreateUC       domain.Create
	DeleteUC       domain.Delete
	FindUC         domain.Find
	FetchUC        domain.Fetch
	UpdateStatusUC domain.UpdateStatus
	UpdateUC       domain.Update
}

func (in Inbound) RegisterTodoServiceServer() {
	he := &httpEndpoint{
		createUC:       in.CreateUC,
		deleteUC:       in.DeleteUC,
		findUC:         in.FindUC,
		fetchUC:        in.FetchUC,
		updateStatusUC: in.UpdateStatusUC,
		updateUC:       in.UpdateUC,
	}

	se := &sseEndpoint{
		codecJSON: in.CodecJSON,
		clients:   make(map[chan []byte]struct{}),
	}

	ge := &grpcEndpoint{
		createUC:       in.CreateUC,
		deleteUC:       in.DeleteUC,
		findUC:         in.FindUC,
		fetchUC:        in.FetchUC,
		updateStatusUC: in.UpdateStatusUC,
		updateUC:       in.UpdateUC,
	}

	gqe := &gqlEndpoint{
		createUC:       in.CreateUC,
		deleteUC:       in.DeleteUC,
		findUC:         in.FindUC,
		fetchUC:        in.FetchUC,
		updateStatusUC: in.UpdateStatusUC,
		updateUC:       in.UpdateUC,
	}

	//
	in.Router.Endpoint(http.MethodGet, "/todos/:id", he.Find)
	in.Router.Endpoint(http.MethodGet, "/todos", he.Fetch)
	in.Router.Endpoint(http.MethodPost, "/todos", he.Create)
	in.Router.Endpoint(http.MethodPut, "/todos/:id", he.Update)
	in.Router.Endpoint(http.MethodPatch, "/todos/:id/status", he.UpdateStatus)
	in.Router.Endpoint(http.MethodDelete, "/todos/:id", he.Delete)

	//
	in.Router.HandleFunc(http.MethodGet, "/events", se.HandleEvent)
	in.Router.HandleFunc(http.MethodGet, "/trigger-event", se.HandleEvent)

	//
	pb.RegisterTodoServiceServer(in.GRPCServer, ge)

	//
	gqlServer := framework.HandlerGQL(ql.NewExecutableSchema(ql.Config{Resolvers: gqe}))
	in.GQLRouter.Handler(http.MethodPost, "/graphql", gqlServer)

	if in.Config.GetBool("feature.flag.graphql.playground") {
		in.GQLRouter.Handler(http.MethodGet, "/graphql/playground",
			playground.Handler("GraphQL playground", "/graphql"))
	}
}
