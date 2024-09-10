package gql

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func RegisterGQLEndpoint(router *httprouter.Router, cfg config.Config, h *Endpoint) {
	exec := ql.NewExecutableSchema(ql.Config{Resolvers: h})
	gqlServer := handler.NewDefaultServer(exec)
	gqlServer.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		log.Printf("%T \n", err)
		log.Println("CTX", ctx)
		log.Println("ERR", err)

		return gqlerror.Errorf("internal server error")
	})

	router.Handler(http.MethodPost, "/graphql", gqlServer)

	if cfg.GetBool("feature.flag.graphql.playground") {
		router.Handler(http.MethodGet, "/graphql/playground",
			playground.Handler("GraphQL playground", "/graphql"),
		)
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
