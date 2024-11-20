package gql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/vektah/gqlparser/v2/ast"
)

type Option func(*Config)

type Config struct {
	introspection bool
}

func WithIntrospection() Option {
	return func(c *Config) {
		c.introspection = true
	}
}

func Handler(es graphql.ExecutableSchema, opts ...Option) *handler.Server {
	srv := handler.New(es)
	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	srv.AddTransport(transportPOST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	if cfg.introspection {
		srv.Use(extension.Introspection{}) // view data schema
	}
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	return srv
}
