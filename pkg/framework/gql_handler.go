package framework

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/vektah/gqlparser/v2/ast"
)

// GQlOption represents a functional option for configuring the GQLConfig.
type GQlOption func(*GQLConfig)

// GQLConfig holds configuration options for a GraphQL server.
type GQLConfig struct {
	introspection bool // Enables or disables schema introspection.
}

// WithIntrospection returns a GQlOption that enables introspection for the GraphQL server.
// Introspection allows clients to query the server's schema.
func WithIntrospection() GQlOption {
	return func(c *GQLConfig) {
		c.introspection = true
	}
}

// HandlerGQL initializes a new GraphQL server handler with the provided executable schema
// and optional configuration options. The handler is configured with default transports
// and extensions, including support for persisted queries and optional introspection.
func HandlerGQL(es graphql.ExecutableSchema, opts ...GQlOption) *handler.Server {
	srv := handler.New(es)
	cfg := &GQLConfig{}

	// Apply provided options to the configuration.
	for _, opt := range opts {
		opt(cfg)
	}

	// Configure the server.
	srv.AddTransport(transportPOST{})                    // Add support for POST requests.
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000)) // Cache query documents.

	if cfg.introspection {
		srv.Use(extension.Introspection{}) // Enable schema introspection if configured.
	}

	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)}) // Enable persisted queries.

	return srv
}
