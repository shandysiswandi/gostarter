package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/framework/gql"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
	"github.com/shandysiswandi/gostarter/pkg/framework/middleware"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
)

func RegisterGQLEndpoint(router *httpserver.Router, h *Endpoint, cfg config.Config, jwte jwt.JWT) {
	exec := ql.NewExecutableSchema(ql.Config{Resolvers: h})
	gqlServer := gql.ServerDefault(exec)
	// gqlServer := handler.New(exec)

	// gqlServer.AddTransport(transport.Websocket{
	// 	KeepAlivePingInterval: 10 * time.Second,
	// })
	// gqlServer.AddTransport(transport.Options{})
	// gqlServer.AddTransport(transport.GET{})
	// gqlServer.AddTransport(transport.POST{
	// 	ResponseHeaders: map[string][]string{
	// 		"Content-Type": {"application/json; charset=utf-8"},
	// 	},
	// })
	// gqlServer.AddTransport(transport.MultipartForm{})

	// gqlServer.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// gqlServer.Use(extension.Introspection{})
	// gqlServer.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New[string](100),
	// })

	// gqlServer.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
	// 	log.Printf("%T \n", err)
	// 	log.Println("CTX", ctx)
	// 	log.Println("ERR", err)

	// 	return gqlerror.Errorf("internal server error")
	// })

	handler := middleware.Chain(gqlServer, middleware.JWT(jwte, "gostarter.access.token"))

	router.Native(http.MethodPost, "/graphql", handler)

	if cfg.GetBool("feature.flag.graphql.playground") {
		router.Native(http.MethodGet, "/graphql/playground",
			playground.Handler("GraphQL playground", "/graphql"))
	}
}

type Endpoint struct {
	ql.Resolver

	FindUC         domain.Find
	FetchUC        domain.Fetch
	CreateUC       domain.Create
	DeleteUC       domain.Delete
	UpdateUC       domain.Update
	UpdateStatusUC domain.UpdateStatus
}

func (e *Endpoint) Mutation() ql.MutationResolver { return &mutation{e} }

func (e *Endpoint) Query() ql.QueryResolver { return &query{e} }
