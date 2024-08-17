//nolint:ireturn // ignore for some reason
package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/config"
)

func RegisterGQLEndpoint(router *httprouter.Router, cfg config.Config, h *Endpoint) {
	exec := ql.NewExecutableSchema(ql.Config{Resolvers: h})
	gqlServer := handler.NewDefaultServer(exec)

	router.Handler(http.MethodPost, "/graphql", gqlServer)

	if cfg.GetBool("feature.flag.graphql.playground") {
		router.Handler(http.MethodGet, "/graphql/playground",
			playground.Handler("GraphQL playground", "/graphql"),
		)
	}
}

type Endpoint struct {
	ql.Resolver

	GetByIDUC       usecase.GetByID
	GetWithFilterUC usecase.GetWithFilter
	CreateUC        usecase.Create
	DeleteUC        usecase.Delete
	UpdateUC        usecase.Update
	UpdateStatusUC  usecase.UpdateStatus
}

func (e *Endpoint) Mutation() ql.MutationResolver { return &mutation{e} }

func (e *Endpoint) Query() ql.QueryResolver { return &query{e} }
